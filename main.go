package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var config Config

type Config struct {
	MonthlyDeposit       int     `yaml:"monthlyDeposit"`
	AlreadyInvested     int     `yaml:"alreadyInvested"`
	MonthlyGrowthPercent float64 `yaml:"monthlyGrowthPercent"`
	DividentPercent     float64 `yaml:"dividentPercent"`
	FreeLifeMoney       int     `yaml:"freeLifeMoney"`
}

func main() {
	err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config file: ", err)
		os.Exit(1)
	}

	calcPlan()
}

func loadConfig() error {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	return nil
}

func calcPlan() {
	var monthlyDeposit int = config.MonthlyDeposit
	var alreadyInvested int = config.AlreadyInvested
	var monthlyGrowthPercent float64 = config.MonthlyGrowthPercent
	var dividentPercent float64 = config.DividentPercent

	var truelyInvested = alreadyInvested
	var moneyOnAcc = int(alreadyInvested)

	var isFreeLife = false

	for i := 18; i <= 80; i++ {
		yearProfit := moneyOnAcc
		yearDividend := 0
		for y := 1; y <= 12; y++ {
			monthProfit := moneyOnAcc
			truelyInvested = truelyInvested + monthlyDeposit                                       // add monthly deposit to truely invested
			moneyOnAcc = moneyOnAcc + monthlyDeposit                                               // add monthly deposit
			moneyOnAcc = int(float64(moneyOnAcc) + (float64(moneyOnAcc/100))*monthlyGrowthPercent) // add monthly growth

			if y%3 == 0 {
				cvartalPercent := dividentPercent / 4
				qurterDividend := int(float64(moneyOnAcc) * cvartalPercent / 100)
				yearDividend = yearDividend + qurterDividend
				moneyOnAcc = int(float64(moneyOnAcc) + float64(qurterDividend)) // add divident
			}
			monthProfit = moneyOnAcc - monthProfit

			if monthProfit > config.FreeLifeMoney && !isFreeLife {
				isFreeLife = true
				fmt.Printf("	Free life at %v years\n", i)
				fmt.Printf("	Month profit:		%v CZK\n", monthProfit)
				fmt.Printf("	Truely invested: 	%v CZK\n", truelyInvested)
				fmt.Printf("	Money growth:		%v CZK\n", moneyOnAcc-truelyInvested)
				fmt.Printf("	Money on account: 	%v CZK\n", moneyOnAcc)
			}
		}

		yearProfit = moneyOnAcc - yearProfit
		fmt.Printf("Y-%v, YP&DEP-%v", i, yearProfit)

		if i%10 == 0 {
			fmt.Printf("---%v---\n", i)
			fmt.Printf("Truely invested: 	%v CZK\n", truelyInvested)
			fmt.Printf("Money growth:		%v CZK\n", moneyOnAcc-truelyInvested)
			fmt.Printf("Money on account: 	%v CZK\n", moneyOnAcc)
		}
	}
}
