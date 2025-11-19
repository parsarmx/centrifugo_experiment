package otp

import (
	"crypto/rand"
	"fmt"
	"math/bits"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"

	"go.uber.org/zap"
)

var pow10Map = map[int]int{
	4: 1000,
	5: 10000,
	6: 100000,
	7: 1000000,
	8: 10000000,
	9: 100000000,
}

type OTPGenerator interface {
	GenerateOTPLowMemoryMode() (otpCode string, err error)
}

type lowMemoryOTPGenerator struct {
	min      int
	max      int
	rangeVal int
	pool     *sync.Pool

	// Add entropy buffer for high RPS scenarios
	randBuf    []byte
	randBufPos uint32
	randBufMu  sync.Mutex

	// Cache of pre-computed strings for common OTP values
	stringCache  []string
	useCacheFlag bool
}

type otpGenerator struct {
	lowMemory lowMemoryOTPGenerator
	logger    *zap.Logger
	OTPLength int
}

func NewDefaultOTPGenerator(OTPLength int, logger *zap.Logger) (OTPGenerator, error) {
	if OTPLength < 4 || OTPLength > 8 {
		logger.Debug("OTP length must be between 4 and 8", zap.Int("OTP_length", OTPLength))
		return nil, fmt.Errorf("OTP length must be between 4 and 8")
	}

	var pool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 8)
		},
	}

	min := pow10Map[OTPLength]
	max := pow10Map[OTPLength+1] - 1

	rangeVal := max - min + 1

	// Pre-allocate a buffer of entropy for high-performance scenarios
	randBuf := make([]byte, 4096) // 4KB of entropy, enough for 512 OTPs
	_, err := rand.Read(randBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize random buffer: %w", err)
	}

	// Pre-compute string cache for common OTP values if OTP length is small enough
	var stringCache []string
	useCacheFlag := false

	// Only use cache for smaller OTP lengths where memory usage is reasonable
	if OTPLength <= 6 {
		useCacheFlag = true
		stringCache = make([]string, rangeVal)
		for i := 0; i < rangeVal; i++ {
			stringCache[i] = strconv.FormatInt(int64(min+i), 10)
		}
	}

	return &otpGenerator{
		logger: logger,
		lowMemory: lowMemoryOTPGenerator{
			min:          min,
			max:          max,
			rangeVal:     rangeVal,
			pool:         &pool,
			randBuf:      randBuf,
			randBufPos:   0,
			stringCache:  stringCache,
			useCacheFlag: useCacheFlag,
		},
		OTPLength: OTPLength,
	}, nil
}

func (d *otpGenerator) GenerateOTPLowMemoryMode() (otpCode string, err error) {

	// Use fast entropy source for high RPS
	randomUint, err := d.getFastRandomUint64()
	if err != nil {
		return "", err
	}

	// Fast modulo using the range - this produces an even distribution
	randomInt := int(fastMod(randomUint, uint64(d.lowMemory.rangeVal)))
	result := d.lowMemory.min + randomInt

	// Use pre-computed string cache for common values if available
	if d.lowMemory.useCacheFlag {
		if randomInt < len(d.lowMemory.stringCache) {
			return d.lowMemory.stringCache[randomInt], nil
		}
	}

	// Format the integer to string with pre-allocated capacity
	return strconv.FormatInt(int64(result), 10), nil
}

func fastMod(n, d uint64) uint64 {
	// Check if d is a power of 2
	if bits.OnesCount64(d) == 1 {
		return n & (d - 1)
	}

	// For small divisors, standard modulo is fine
	return n % d
}

// getFastRandomUint64 gets random bytes from the buffer with minimal locking
func (d *otpGenerator) getFastRandomUint64() (uint64, error) {
	// Use buffered entropy first - this is significantly faster for high RPS
	pos := atomic.AddUint32(&d.lowMemory.randBufPos, 8) - 8

	// If we're near the end of the buffer, refill it
	if pos >= uint32(len(d.lowMemory.randBuf)-8) {
		d.lowMemory.randBufMu.Lock()
		// Double-check if another goroutine already refilled the buffer
		if atomic.LoadUint32(&d.lowMemory.randBufPos) >= uint32(len(d.lowMemory.randBuf)-8) {
			_, err := rand.Read(d.lowMemory.randBuf)
			if err != nil {
				d.lowMemory.randBufMu.Unlock()
				return 0, fmt.Errorf("failed to refill random buffer: %w", err)
			}
			atomic.StoreUint32(&d.lowMemory.randBufPos, 0)
			pos = 0
		} else {
			// Another goroutine already refilled, get new position
			pos = atomic.AddUint32(&d.lowMemory.randBufPos, 8) - 8
		}
		d.lowMemory.randBufMu.Unlock()
	}

	// Read 8 bytes directly from the buffer at the current position
	return *(*uint64)(unsafe.Pointer(&d.lowMemory.randBuf[pos])), nil
}
