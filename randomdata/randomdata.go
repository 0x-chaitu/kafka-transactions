package randomdata

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var locations = []string{
	"New York, NY",
	"Los Angeles, CA",
	"Chicago, IL",
	"Houston, TX",
	"Phoenix, AZ",
	"Philadelphia, PA",
	"San Antonio, TX",
	"San Diego, CA",
	"Dallas, TX",
	"San Jose, CA",
	"Austin, TX",
	"Jacksonville, FL",
	"Fort Worth, TX",
	"Columbus, OH",
	"Charlotte, NC",
	"San Francisco, CA",
	"Indianapolis, IN",
	"Seattle, WA",
	"Denver, CO",
	"Washington, DC",
}

func AccountNumber() int {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	r.Seed(time.Now().UnixNano())
	return r.Intn(999999999-111111111+1) + 111111111
}

func TransactionId() int {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	r.Seed(time.Now().UnixNano())
	return r.Intn(999999999-111111111+1) + 111111111
}

func TransactionAmount(minAmount, maxAmount float32) float32 {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	randomAmount := r.Float32()*(maxAmount-minAmount) + minAmount
	formattedAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", randomAmount), 32)
	return float32(formattedAmount)
}

func TransactionTime() time.Time {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	randomDuration := time.Duration(r.Intn(86400) * int(time.Second))
	randomTime := time.Now().Add(-randomDuration)
	return randomTime
}

func Location() string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return locations[r.Intn(len(locations))]
}
