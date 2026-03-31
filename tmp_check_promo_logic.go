package main

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/model"
	"fmt"
	"time"
	_ "time/tzdata"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.ConnectDB()

	var promo model.Promo
	err := config.DB.First(&promo, 5).Error
	if err != nil {
		fmt.Println("Error taking promo 5:", err)
		os.Exit(1)
	}

	fmt.Printf("Promo ID: %d\n", promo.ID)
	fmt.Printf("Name: %s\n", promo.Name)
	fmt.Printf("Status: %s\n", promo.Status)
	
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("Error loading location, fallback to FixedZone:", err)
		loc = time.FixedZone("WIB", 7*3600)
	}
	
	now := time.Now().In(loc)
	fmt.Printf("\nNow (Asia/Jakarta): %v\n", now)

	// NEW LOGIC
	fmt.Println("\n--- Validating (Tz Aware Logic) ---")
	
	startTime := promo.StartTime
	if startTime == "" { startTime = "00:00" }
	endTime := promo.EndTime
	if endTime == "" { endTime = "23:59" }

	startStr := fmt.Sprintf("%s %s", promo.StartDate.Format("2006-01-02"), startTime)
	endStr := fmt.Sprintf("%s %s", promo.EndDate.Format("2006-01-02"), endTime)

	startDT, _ := time.ParseInLocation("2006-01-02 15:04", startStr, loc)
	endDT, _ := time.ParseInLocation("2006-01-02 15:04", endStr, loc)

	// Handle StartTime == EndTime for campaign bounds
	if startTime == endTime && startTime != "00:00" {
		endDT = endDT.Add(24 * time.Hour)
		fmt.Println("StartTime == EndTime detected, extended endDT by 24h")
	}

	fmt.Printf("Parsed StartDT: %v\n", startDT)
	fmt.Printf("Parsed EndDT: %v\n", endDT)

	validDT := true
	if now.Before(startDT) || now.After(endDT) {
		fmt.Println("Failed Absolute DateTime check")
		validDT = false
	}
	fmt.Printf("DateTime check: %v\n", validDT)
}
