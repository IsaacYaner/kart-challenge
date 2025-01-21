package service

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	SHARD_COUNT = 256 // Increase shard count to reduce lock contention
	BUFFER_SIZE = 4 * 1024 * 1024
)

type CouponShard struct {
	coupons      map[uint64]uint8  // hash -> file bits (001, 010, 100 for each file)
	validCoupons map[uint64]string // hash -> original coupon (only for valid ones)
	mu           sync.RWMutex
}

type CouponService struct {
	shards []*CouponShard
}

func NewCouponShard() *CouponShard {
	return &CouponShard{
		coupons:      make(map[uint64]uint8),
		validCoupons: make(map[uint64]string),
	}
}

func NewCouponService() *CouponService {
	cs := &CouponService{
		shards: make([]*CouponShard, SHARD_COUNT),
	}
	for i := 0; i < SHARD_COUNT; i++ {
		cs.shards[i] = NewCouponShard()
	}

	if err := cs.initializeCoupons(); err != nil {
		log.Printf("Error initializing coupons: %v", err)
	}
	return cs
}

// getShard returns the appropriate shard for a given coupon
func (s *CouponService) getShard(hash uint64) *CouponShard {
	return s.shards[hash%SHARD_COUNT]
}

func hashString(s string) uint64 {
	h := uint64(0)
	for i := 0; i < len(s); i++ {
		h = h*31 + uint64(s[i])
	}
	return h
}

func (s *CouponService) insert(coupon string, fileIndex uint8) {
	hash := hashString(coupon)
	shard := s.getShard(hash)

	shard.mu.Lock()
	fileBit := uint8(1 << fileIndex) // 001, 010, or 100 for file1, file2, file3
	currentBits := shard.coupons[hash]
	newBits := currentBits | fileBit // Combine bits with OR operation

	if currentBits != newBits { // Only update if it's a new file
		shard.coupons[hash] = newBits
		// Store original coupon when it appears in 2+ files
		if countBits(newBits) > 1 && countBits(currentBits) == 1 {
			shard.validCoupons[hash] = coupon
		}
	}
	shard.mu.Unlock()
}

// countBits counts the number of 1s in a byte
func countBits(b uint8) int {
	count := 0
	for b != 0 {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

func (s *CouponService) generateValidCoupons() error {
	couponFiles := []string{
		"internal/data/coupons/couponbase1.txt",
		"internal/data/coupons/couponbase2.txt",
		"internal/data/coupons/couponbase3.txt",
	}

	var wg sync.WaitGroup
	errorChan := make(chan error, len(couponFiles))

	for fileIndex, filename := range couponFiles {
		wg.Add(1)
		go func(idx uint8, fname string) {
			defer wg.Done()

			file, err := os.Open(fname)
			if err != nil {
				errorChan <- err
				return
			}
			defer file.Close()

			reader := bufio.NewReaderSize(file, BUFFER_SIZE)
			scanner := bufio.NewScanner(reader)
			scanner.Buffer(make([]byte, BUFFER_SIZE), BUFFER_SIZE)

			var processed int
			for scanner.Scan() {
				if coupon := scanner.Text(); coupon != "" {
					s.insert(coupon, idx)
					processed++
					if processed%1000000 == 0 {
						log.Printf("Processed %d coupons from %s", processed, fname)
					}
				}
			}

			if err := scanner.Err(); err != nil {
				log.Printf("Error scanning file %s: %v", fname, err)
				errorChan <- err
				return
			}

			log.Printf("Completed processing %d coupons from %s", processed, fname)
		}(uint8(fileIndex), filename)
	}

	wg.Wait()
	close(errorChan)

	// Check for errors
	for err := range errorChan {
		if err != nil {
			return errors.New("error processing coupon files")
		}
	}

	return s.writeValidCoupons()
}

func (s *CouponService) writeValidCoupons() error {
	file, err := os.Create("internal/data/coupons/valid_coupons.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Process one shard at a time
	for _, shard := range s.shards {
		shard.mu.RLock()
		for hash, bits := range shard.coupons {
			if countBits(bits) > 1 {
				if coupon, exists := shard.validCoupons[hash]; exists {
					// Format: coupon [file1,file2,file3]
					files := getFileList(bits)
					if _, err := writer.WriteString(fmt.Sprintf("%s [%s]\n", coupon, files)); err != nil {
						shard.mu.RUnlock()
						return err
					}
				}
			}
		}
		shard.mu.RUnlock()
	}
	return nil
}

// getFileList converts bit field to file list string
func getFileList(bits uint8) string {
	var files []string
	fileNames := []string{"file1", "file2", "file3"}
	for i := uint8(0); i < 3; i++ {
		if bits&(1<<i) != 0 {
			files = append(files, fileNames[i])
		}
	}
	return strings.Join(files, ",")
}

func (s *CouponService) IsValidCoupon(code string) bool {
	hash := hashString(code)
	shard := s.getShard(hash)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	count := shard.coupons[hash]
	return count > 1
}

func readCouponsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var coupons []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		coupon := scanner.Text()
		if coupon != "" {
			coupons = append(coupons, coupon)
		}
	}
	return coupons, scanner.Err()
}

func (s *CouponService) initializeCoupons() error {
	validCouponFile := "internal/data/coupons/valid_coupons.txt"

	if _, err := os.Stat(validCouponFile); os.IsNotExist(err) {
		if err := s.generateValidCoupons(); err != nil {
			log.Printf("Failed to generate valid coupons: %v", err)
			return errors.New("failed to generate valid coupons")
		}
	}

	return s.loadValidCoupons(validCouponFile)
}

func (s *CouponService) loadValidCoupons(filename string) error {
	coupons, err := readCouponsFromFile(filename)
	if err != nil {
		return err
	}

	for _, coupon := range coupons {
		s.insert(coupon, 0)
	}
	return nil
}
