package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
)

type TreeStructure struct {
	UserID              int
	ParentID            *TreeStructure
	Children            []*TreeStructure
	Joining_package_fee float64
	BV                  float64
	Position            string
	Levels              int
	LeftMember          *TreeStructure
	RightMember         *TreeStructure
	SponsorBonus        float64
	BinaryBonus         float64
	CarryForward        float64
	MatchingBonus       float64
	CappingValue        float64
	Cycle               int
	FlushOut            float64
	LeftCarry           float64
	RightCarry          float64
	LeftDownlineSale    float64
	RightDownlineSale   float64
	CycleProfit         float64
	PoolBonus           float64
}

type TreeStructureJSON struct {
	UserID              int     `json:"user_id"`
	Position            string  `json:"position"`
	ParentID            int     `json:"parent_id"`
	Joining_package_fee float64 `json:"joining_package_fee"`
	BV                  float64 `json:"bv"`
	LeftMember          int     `json:"left_member"`
	RightMember         int     `json:"right_member"`
	Levels              int     `json:"levels"`
	SponsorBonus        float64 `json:"sponsor_bonus"`
	BinaryBonus         float64 `json:"binary_bonus"`
	MatchingBonus       float64 `json:"matching_bonus"`
	CarryForward        float64 `json:"carry_forward"`
	Cycle               int     `json:"cycle"`
	FlushOut            float64 `json:"flush_out"`
	LeftCarry           float64 `json:"left_carry"`
	RightCarry          float64 `json:"right_carry"`
	LeftDownlineSale    float64 `json:"LeftDownlineSale"`
	RightDownlineSale   float64 `json:"RightDownlineSale"`
	CycleProfit         float64 `json:"cycle_profit"`
	PoolBonus           float64 `json:"pool_bonus"`
}

// func convertToJSONStructure(nodes []*TreeStructure) []TreeStructureJSON {
// 	var jsonNodes []TreeStructureJSON
// 	for _, node := range nodes {
// 		var parentID int
// 		var leftMemberID, rightMemberID int
// 		if node.ParentID != nil {
// 			parentID = node.ParentID.UserID
// 		} else {
// 			parentID = 0
// 		}
// 		if node.LeftMember != nil {
// 			leftMemberID = node.LeftMember.UserID
// 		} else {
// 			leftMemberID = 0
// 		}

// 		if node.RightMember != nil {
// 			rightMemberID = node.RightMember.UserID
// 		} else {
// 			rightMemberID = 0
// 		}

// 		jsonNodes = append(jsonNodes, TreeStructureJSON{
// 			UserID:              node.UserID,
// 			Position:            node.Position,
// 			Joining_package_fee: node.Joining_package_fee,
// 			BV:                  node.BV,
// 			LeftMember:          leftMemberID,
// 			RightMember:         rightMemberID,
// 			Levels:              node.Levels,
// 			ParentID:            parentID,
// 			SponsorBonus:        node.SponsorBonus,
// 			BinaryBonus:         node.BinaryBonus,
// 			PoolBonus:           node.PoolBonus,
// 			MatchingBonus:       node.MatchingBonus,
// 			CycleProfit:         node.CycleProfit,
// 			CarryForward:        node.CarryForward,
// 			Cycle:               node.Cycle,
// 			FlushOut:            node.FlushOut,
// 			LeftCarry:           node.LeftCarry,
// 			RightCarry:          node.RightCarry,
// 			LeftDownlineSale:    node.LeftDownlineSale,
// 			RightDownlineSale:   node.RightDownlineSale,
// 		})
// 	}
// 	return jsonNodes
// }

var treeNodes []*TreeStructure

func buildTree(root *TreeStructure, numMembers int, startID int) []*TreeStructure {
	if numMembers <= 0 {
		return nil
	}

	treeNodes := []*TreeStructure{root}
	queue := []*TreeStructure{root}
	currentID := startID
	nodesToAdd := numMembers
	for len(queue) > 0 && nodesToAdd > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.LeftMember == nil && nodesToAdd > 0 {
			leftChild := &TreeStructure{
				UserID:   currentID,
				ParentID: node,
				Levels:   node.Levels + 1,
				Position: "Left",
			}
			node.LeftMember = leftChild
			queue = append(queue, leftChild)
			treeNodes = append(treeNodes, leftChild)
			currentID++
			nodesToAdd--

			if nodesToAdd > 0 {
				rightChild := &TreeStructure{
					UserID:   currentID,
					ParentID: node,
					Levels:   node.Levels + 1,
					Position: "Right",
				}
				node.RightMember = rightChild
				queue = append(queue, rightChild)
				treeNodes = append(treeNodes, rightChild)
				currentID++
				nodesToAdd--
			}
		} else {
			rightChild := &TreeStructure{
				UserID:   currentID,
				ParentID: node,
				Levels:   node.Levels + 1,
				Position: "Right",
			}
			node.RightMember = rightChild
			queue = append(queue, rightChild)
			treeNodes = append(treeNodes, rightChild)
			currentID++
			nodesToAdd--
		}
	}
	return treeNodes
}

var all_data [][]*TreeStructure

func AssignJoiningFee(nodes []*TreeStructure, joining_package_fee []float64, b_volume []float64, product_quantity []int) {
	if len(nodes) <= 1 {
		return
	}
	allExceptFirst := nodes[1:]
	var jpf, bv float64
	iterant := 0
	var assig_list, bv_assig_list []float64

	for count, i := range product_quantity {
		for j := 0; j < i; j = j + 1 {
			assig_list = append(assig_list, joining_package_fee[count])
			bv_assig_list = append(bv_assig_list, b_volume[count])
		}
	}
	for _, node := range allExceptFirst {
		if iterant < len(assig_list) {
			jpf = assig_list[iterant]
			bv = bv_assig_list[iterant]
		} else {
			iterant = 0
			jpf = joining_package_fee[iterant]
			bv = b_volume[iterant]
		}
		node.Joining_package_fee = jpf
		node.BV = bv
		iterant = iterant + 1
	}
}

func AllocateMembers(numMembers int, product_quantity []int, startID int, result []*TreeStructure, joining_package_fee []float64, b_volume []float64) ([][]int, int, [][]*TreeStructure, map[int]float64) {
	currentID := startID
	remaining := numMembers
	var allCycles [][]int
	totalJoiningFee := make(map[int]float64)
	cycleCount := 1
	var all_data [][]*TreeStructure

	for remaining > 0 {
		var cycle [][]int
		for i, qty := range product_quantity {
			cycleID := []int{}

			assignQty := int(math.Min(float64(qty), float64(remaining)))
			for j := 0; j < assignQty; j++ {
				cycleID = append(cycleID, currentID)
				currentID++
				remaining--
			}
			cycle = append(cycle, cycleID)

			totalJoiningFee[i] += float64(assignQty) * joining_package_fee[i]
			if remaining <= 0 {
				break
			}
		}

		var temp []*TreeStructure
		for _, node := range result {
			if node.UserID < currentID {
				temp = append(temp, node)
			}
		}
		AssignJoiningFee(temp, joining_package_fee, b_volume, product_quantity)
		all_data = append(all_data, temp)
		cycleCount++

		for _, lst := range cycle {
			if len(lst) > 0 {
				allCycles = append(allCycles, lst)
			}
		}
	}
	return allCycles, cycleCount, all_data, totalJoiningFee
}

func CalculateSponsorBonus(allData [][]*TreeStructure, sponsorBonusPercent float64, sponsor_bonus_type1 string, sponsor_bonus_type2 string, cappingAmount float64, cappingScope, bonusOption string) (float64, map[int]float64) {
	var TSB float64

	totalBonus := make(map[int]float64)

	calculateBonus := func(member *TreeStructure, valueSelector func(*TreeStructure) float64) float64 {
		var rightBonus, leftBonus float64
		// var selected string

		if member.RightMember != nil {
			rightBonus = valueSelector(member.RightMember)
		}
		if member.LeftMember != nil {
			leftBonus = valueSelector(member.LeftMember)
		}

		// if strings.Contains(sponsor_bonus_type1, "percent") || strings.Contains(sponsor_bonus_type2, "percent") {
		// 	// if strings.Contains(sponsor_bonus_type1, "percent") {
		// 	// 	selected = sponsor_bonus_type1
		// 	// } else {
		// 	// 	selected = sponsor_bonus_type2
		// 	// }
		// 	if member.RightMember != nil {
		// 		rightBonus = valueSelector(member.RightMember) * sponsorBonusPercent / 100
		// 	}
		// 	if member.LeftMember != nil {
		// 		leftBonus = valueSelector(member.LeftMember) * sponsorBonusPercent / 100
		// 	}
		// } else if strings.Contains(sponsor_bonus_type1, "usd") || strings.Contains(sponsor_bonus_type2, "usd") {
		// 	// if strings.Contains(sponsor_bonus_type1, "usd") {
		// 	// 	selected = sponsor_bonus_type1
		// 	// } else {
		// 	// 	selected = sponsor_bonus_type2
		// 	// }
		// 	if member.RightMember != nil {
		// 		member.SponsorBonus = sponsorBonusPercent
		// 		rightBonus = member.SponsorBonus
		// 	}
		// 	if member.LeftMember != nil {
		// 		member.SponsorBonus = sponsorBonusPercent
		// 		leftBonus = member.SponsorBonus
		// 	}
		// }
		if strings.Contains(sponsor_bonus_type1, "usd") {
			if member.RightMember != nil {
				member.SponsorBonus = sponsorBonusPercent
				rightBonus = member.SponsorBonus
			}
			if member.LeftMember != nil {
				member.SponsorBonus = sponsorBonusPercent
				leftBonus = member.SponsorBonus
			}
		} else {
			if member.RightMember != nil {
				rightBonus = valueSelector(member.RightMember) * sponsorBonusPercent / 100
			}
			if member.LeftMember != nil {
				leftBonus = valueSelector(member.LeftMember) * sponsorBonusPercent / 100
			}
		}
		
		TSB := rightBonus + leftBonus
		if strings.Contains(cappingScope, "sponsor") && TSB > cappingAmount {
			TSB = cappingAmount
		}
		return TSB
	}
	for levelIndex, members := range allData {
		for _, member := range members {
			var bonus float64
			if strings.Contains(bonusOption, "PRICE") {
				bonus = calculateBonus(member, func(m *TreeStructure) float64 { return m.Joining_package_fee })
			} else if strings.Contains(bonusOption, "BV") {
				bonus = calculateBonus(member, func(m *TreeStructure) float64 { return m.BV })
			}
			member.SponsorBonus = bonus
			if member.UserID != 1 {
				TSB += bonus
				totalBonus[levelIndex] += bonus
			}
		}
	}
	return TSB, totalBonus
}

func Traverse(node *TreeStructure, limit int, bonus_option string) float64 {
	if node == nil || node.UserID > limit {
		return 0
	}
	var currentSales float64
	if strings.Contains(bonus_option, "PRICE") {
		currentSales = node.Joining_package_fee
	}
	if strings.Contains(bonus_option, "BV") {
		currentSales = node.BV
	}
	leftSales := Traverse(node.LeftMember, limit, bonus_option)
	rightSales := Traverse(node.RightMember, limit, bonus_option)
	return currentSales + leftSales + rightSales
}

func BinaryWithRatio(allData [][]*TreeStructure, bonusOption string, binaryRatio string, ratioAmount int, cappingScope string, cappingAmount float64, cycleCount int, binaryPercentage float64, binary_bonus_type string) (float64, map[int]float64) {
	var TBB, nodeBonus float64
	totalBonus := make(map[int]float64)
	b1 := 5.0
	b2 := 10.0
	b3 := 15.0
	for i, nodeList := range allData {
		for _, node := range nodeList {
			if i != 0 && node.UserID < len(allData[i-1]) {
				node.LeftCarry = allData[i-1][node.UserID-1].LeftCarry
				node.RightCarry = allData[i-1][node.UserID-1].RightCarry
			}
			node.LeftDownlineSale = float64(node.LeftCarry)
			node.RightDownlineSale = float64(node.RightCarry)
			var leftDownline, rightDownline float64
			if node.LeftMember != nil {
				leftDownline = Traverse(node.LeftMember, len(nodeList), bonusOption)
			}
			if node.RightMember != nil {
				rightDownline = Traverse(node.RightMember, len(nodeList), bonusOption)
			}
			node.LeftDownlineSale = float64(leftDownline) + float64(node.LeftCarry)
			node.RightDownlineSale = float64(rightDownline) + float64(node.RightCarry)

			var left, right, noOfPairs int
			var carryLeft, carryRight float64
			switch binaryRatio {
			case "1":
				left, right = ratioAmount*1, ratioAmount*1
			case "2":
				left, right = ratioAmount*1, ratioAmount*2
			case "3":
				left, right = ratioAmount*2, ratioAmount*1
			}

			leftTemp := int(math.Floor(float64(node.LeftDownlineSale) / float64(left)))
			rightTemp := int(math.Floor(float64(node.RightDownlineSale) / float64(right)))
			noOfPairs = int(math.Min(float64(leftTemp), float64(rightTemp)))
			carryLeft = node.LeftDownlineSale - float64(noOfPairs*left)
			carryRight = node.RightDownlineSale - float64(noOfPairs*right)

			node.LeftCarry = carryLeft
			node.RightCarry = carryRight
			minimumVal := int(math.Min(float64(noOfPairs*left), float64(noOfPairs*right)))

			var bonusPerc float64
			switch {
			case noOfPairs > 0 && noOfPairs <= 5 && binaryPercentage == -1:
				bonusPerc = b1
			case noOfPairs > 5 && noOfPairs <= 10 && binaryPercentage == -1:
				bonusPerc = b2
			case noOfPairs > 10 && binaryPercentage == -1:
				bonusPerc = b3
			default:
				bonusPerc = binaryPercentage
			}
			if strings.Contains(binary_bonus_type, "percent") {
				nodeBonus = (float64(minimumVal) * bonusPerc) / 100
			} else if strings.Contains(binary_bonus_type, "usd") {
				bonusPerc = (float64(minimumVal) * bonusPerc)
			}
			if strings.Contains(cappingScope, "binary") {
				if nodeBonus > cappingAmount {
					node.FlushOut = float64(nodeBonus) - cappingAmount
					nodeBonus = cappingAmount
				}
			}
			node.BinaryBonus = float64(nodeBonus)
			if node.UserID != 1 {
				TBB += node.BinaryBonus
				totalBonus[i] += node.BinaryBonus
			}
		}
	}
	return TBB, totalBonus
}

func CalculateMatchingBonus(allData [][]*TreeStructure, matchingPercentages []float64, cappingAmount float64, cappingScope string) (float64, map[int]float64) {
	var TMB float64
	totalBonus := make(map[int]float64)
	for _, members := range allData {
		for _, member := range members {
			iterant := 0
			if member.ParentID == nil {
				continue
			}
			parent := member.ParentID
			ApplyMatchingBonus(member, parent, matchingPercentages, iterant, cappingAmount, cappingScope)

		}
	}
	for i, list := range allData {
		for _, member := range list {
			if strings.Contains(cappingScope, "matching") && member.MatchingBonus > cappingAmount {
				member.MatchingBonus = cappingAmount
			}
			if member.UserID != 1 {
				TMB += member.MatchingBonus
				totalBonus[i] += member.MatchingBonus
			}
		}
	}
	return TMB, totalBonus
}

func ApplyMatchingBonus(member *TreeStructure, parent *TreeStructure, matchingPercentages []float64, iterant int, cappingAmount float64, cappingScope string) {
	if iterant >= len(matchingPercentages) || parent == nil {
		return
	}
	matching_bonus := parent.MatchingBonus + (member.BinaryBonus * matchingPercentages[iterant] / 100)

	parent.MatchingBonus = matching_bonus
	iterant = iterant + 1
	parent = parent.ParentID
	ApplyMatchingBonus(member, parent, matchingPercentages, iterant, cappingAmount, cappingScope)
}

func FindProfitTOCompany(all_Data [][]*TreeStructure, ExpenseMembers float64, totalJoiningCycle map[int]float64) (float64, float64, map[int]float64, map[int]float64, map[int]float64) {
	totalProfit := make(map[int]float64)
	totalExpense := make(map[int]float64)
	totalRevenue := make(map[int]float64)
	var cycleProfits, Revenue, Expense float64
	for i, list := range all_Data {
		var expense, cycleRevenue float64
		for _, member := range list {
			if member.UserID > 1 {
				cycleRevenue = cycleRevenue + member.Joining_package_fee
				expense = expense + ExpenseMembers + member.BinaryBonus + member.MatchingBonus + member.SponsorBonus
			}
		}
		profit := cycleRevenue - expense
		cycleProfits = cycleProfits + profit
		Expense = Expense + expense
		Revenue = Revenue + cycleRevenue
		totalRevenue[i] = Revenue
		totalProfit[i] = profit
		totalExpense[i] = Expense
	}
	return cycleProfits, Revenue, totalProfit, totalExpense, totalRevenue
}

func CalculatePoolBonus(all_Data [][]*TreeStructure, totalProfit map[int]float64, pool_perc float64, dist_no int) (float64, map[int]float64) {
	var pool_amount float64
	PoolAmount := make(map[int]float64)

	for i := range totalProfit {
		PoolAmount[i] = totalProfit[i] * pool_perc / 100
		totalProfit[i] = totalProfit[i] - PoolAmount[i]
		pool_amount = pool_amount + PoolAmount[i]
	}
	return pool_amount, PoolAmount
}

// UNILEVEL
func buildUnilevelTree(root *TreeStructure, numMembers int, startID int, numChild int) []*TreeStructure {
	if numMembers <= 0 {
		return nil
	}
	treeNodes := []*TreeStructure{root}
	queue := []*TreeStructure{root}
	currentID := startID
	nodesToAdd := numMembers

	for len(queue) > 0 && nodesToAdd > 0 {
		node := queue[0]
		queue = queue[1:]

		for i := 0; i < numChild && nodesToAdd > 0; i++ {
			child := &TreeStructure{
				UserID:   currentID,
				ParentID: node,
				Levels:   node.Levels + 1,
			}
			node.Children = append(node.Children, child)
			queue = append(queue, child)
			treeNodes = append(treeNodes, child)
			currentID++
			nodesToAdd--
		}
	}
	return treeNodes
}

func AssignUnilevelJoiningFee(nodes []*TreeStructure, joining_package_fee []float64, product_quantity []int) {
	if len(nodes) <= 1 {
		return
	}
	allExceptFirst := nodes[1:]
	var jpf float64
	iterant := 0
	var assig_list []float64

	for count, i := range product_quantity {
		for j := 0; j < i; j = j + 1 {
			assig_list = append(assig_list, joining_package_fee[count])
		}
	}
	for _, node := range allExceptFirst {
		if iterant < len(assig_list) {
			jpf = assig_list[iterant]
		} else {
			iterant = 0
			jpf = joining_package_fee[iterant]
		}
		node.Joining_package_fee = jpf
		iterant = iterant + 1
	}
}

func AllocateUnilevelMembers(numMembers int, product_quantity []int, startID int, result []*TreeStructure, joining_package_fee []float64) ([][]int, int, [][]*TreeStructure, map[int]float64) {
	currentID := startID
	remaining := numMembers
	var allCycles [][]int
	totalJoiningFee := make(map[int]float64)
	cycleCount := 1
	var all_data [][]*TreeStructure

	for remaining > 0 {
		var cycle [][]int
		for i, qty := range product_quantity {
			cycleID := []int{}

			assignQty := int(math.Min(float64(qty), float64(remaining)))
			for j := 0; j < assignQty; j++ {
				cycleID = append(cycleID, currentID)
				currentID++
				remaining--
			}
			cycle = append(cycle, cycleID)

			totalJoiningFee[i] += float64(assignQty) * joining_package_fee[i]

			if remaining <= 0 {
				break
			}
		}

		var temp []*TreeStructure
		for _, node := range result {
			if node.UserID < currentID {
				temp = append(temp, node)
			}
		}
		AssignUnilevelJoiningFee(temp, joining_package_fee, product_quantity)
		all_data = append(all_data, temp)
		cycleCount++

		for _, lst := range cycle {
			if len(lst) > 0 {
				allCycles = append(allCycles, lst)
			}
		}
	}
	return allCycles, cycleCount, all_data, totalJoiningFee
}

func CalculateUnilevelSponsorBonus(allData [][]*TreeStructure, sponsorBonusPercent float64, uni_sponsor_bonus_type string, cappingAmount float64, cappingScope string) (float64, map[int]float64) {
	var TSB, sponsorBonus float64
	totalBonus := make(map[int]float64)
	for _, list := range allData {
		for _, member := range list {
			if member.ParentID != nil {
				if strings.Contains(uni_sponsor_bonus_type, "percent") {
					sponsorBonus = member.ParentID.SponsorBonus + member.Joining_package_fee*sponsorBonusPercent/100
				} else if strings.Contains(uni_sponsor_bonus_type, "usd") {
					sponsorBonus = member.ParentID.SponsorBonus + sponsorBonusPercent
				}
				if strings.Contains(cappingScope, "sponsor") && sponsorBonus > cappingAmount {
					member.ParentID.SponsorBonus = cappingAmount
				} else {
					member.ParentID.SponsorBonus = sponsorBonus
				}
				member.ParentID.SponsorBonus = sponsorBonus
			}
		}
	}

	for i, list := range allData {
		for _, member := range list {

			if member.UserID != 1 {
				TSB += member.SponsorBonus
				totalBonus[i] += member.SponsorBonus
			}
		}
	}
	return TSB, totalBonus
}

func CalculateUnilevelMatchingBonus(allData [][]*TreeStructure, matchingPercentages []float64, cappingAmount float64, cappingScope string) (float64, map[int]float64) {
	var TMB float64
	totalBonus := make(map[int]float64)
	for _, members := range allData {
		for _, member := range members {
			iterant := 0
			if member.ParentID == nil {
				continue
			}
			parent := member.ParentID
			ApplyUnilevelMatchingBonus(member, parent, matchingPercentages, iterant, cappingAmount, cappingScope)

		}
	}
	for i, list := range allData {
		for _, member := range list {
			if strings.Contains(cappingScope, "matching") && member.MatchingBonus > cappingAmount {
				member.MatchingBonus = cappingAmount
			}

			if member.UserID != 1 {
				TMB += member.MatchingBonus
				totalBonus[i] += member.MatchingBonus
			}
		}
	}
	return TMB, totalBonus
}

func ApplyUnilevelMatchingBonus(member *TreeStructure, parent *TreeStructure, matchingPercentages []float64, iterant int, cappingAmount float64, cappingScope string) {
	if iterant >= len(matchingPercentages) || parent == nil {
		return
	}
	matching_bonus := parent.MatchingBonus + (member.SponsorBonus * matchingPercentages[iterant] / 100)
	parent.MatchingBonus = matching_bonus
	iterant = iterant + 1
	parent = parent.ParentID
	ApplyUnilevelMatchingBonus(member, parent, matchingPercentages, iterant, cappingAmount, cappingScope)
}

func FindUnilevelProfitToCompany(all_Data [][]*TreeStructure, ExpenseMembers float64, totalJoiningCycle map[int]float64, pool_perc float64) (float64, map[int]float64, float64, map[int]float64) {
	totalProfit := make(map[int]float64)
	var cycleProfits float64
	var pool_amount float64
	PoolAmount := make(map[int]float64)
	for i, list := range all_Data {
		var expense, cycleRevenue float64
		for _, member := range list {
			if member.UserID != 1 {
				cycleRevenue = cycleRevenue + member.Joining_package_fee
				expense = expense + ExpenseMembers + member.BinaryBonus + member.MatchingBonus + member.SponsorBonus + member.PoolBonus
			}
		}
		profit := cycleRevenue - expense
		expense = expense + (profit * (pool_perc / 100))

		PoolAmount[i] = profit * pool_perc / 100

		profit = cycleRevenue - expense
		cycleProfits = cycleProfits + profit
		totalProfit[i] = profit

		pool_amount = pool_amount + PoolAmount[i]
	}
	return cycleProfits, totalProfit, pool_amount, PoolAmount
}

// func CalculateUnilevelPoolBonus(all_Data [][]*TreeStructure, totalProfit map[int]float64, pool_perc float64, dist_no int) (float64, map[int]float64) {
// 	var pool_amount float64
// 	PoolAmount := make(map[int]float64)

// 	for i := range totalProfit {
// 		PoolAmount[i] = totalProfit[i] * pool_perc / 100
// 		totalProfit[i] = totalProfit[i] - PoolAmount[i]
// 		pool_amount = pool_amount + PoolAmount[i]
// 	}
// 	return pool_amount, PoolAmount
// }

func buildMatrixTree(root *TreeStructure, numMembers int, startID int, numChild int) []*TreeStructure {
	if numMembers <= 0 {
		return nil
	}
	treeNodes := []*TreeStructure{root}
	queue := []*TreeStructure{root}
	currentID := startID
	nodesToAdd := numMembers

	for len(queue) > 0 && nodesToAdd > 0 {
		node := queue[0]
		queue = queue[1:]

		for i := 0; i < numChild && nodesToAdd > 0; i++ {
			child := &TreeStructure{
				UserID:   currentID,
				ParentID: node,
				Levels:   node.Levels + 1,
			}
			node.Children = append(node.Children, child)
			queue = append(queue, child)
			treeNodes = append(treeNodes, child)
			currentID++
			nodesToAdd--
		}
	}
	return treeNodes
}

func AssignMatrixJoiningFee(nodes []*TreeStructure, joining_package_fee []float64, product_quantity []int) {
	if len(nodes) <= 1 {
		return
	}
	allExceptFirst := nodes[1:]
	var jpf float64
	iterant := 0
	var assig_list []float64

	for count, i := range product_quantity {
		for j := 0; j < i; j = j + 1 {
			assig_list = append(assig_list, joining_package_fee[count])
		}
	}
	for _, node := range allExceptFirst {
		if iterant < len(assig_list) {
			jpf = assig_list[iterant]
		} else {
			iterant = 0
			jpf = joining_package_fee[iterant]
		}
		node.Joining_package_fee = jpf
		iterant = iterant + 1
	}
}

func AllocateMatrixMembers(numMembers int, product_quantity []int, startID int, result []*TreeStructure, joining_package_fee []float64) ([][]int, int, [][]*TreeStructure, map[int]float64) {
	currentID := startID
	remaining := numMembers
	var allCycles [][]int
	totalJoiningFee := make(map[int]float64)
	cycleCount := 1
	var all_data [][]*TreeStructure

	for remaining > 0 {
		var cycle [][]int
		for i, qty := range product_quantity {
			cycleID := []int{}

			assignQty := int(math.Min(float64(qty), float64(remaining)))
			for j := 0; j < assignQty; j++ {
				cycleID = append(cycleID, currentID)
				currentID++
				remaining--
			}
			cycle = append(cycle, cycleID)

			totalJoiningFee[i] += float64(assignQty) * joining_package_fee[i]

			if remaining <= 0 {
				break
			}
		}

		var temp []*TreeStructure
		for _, node := range result {
			if node.UserID < currentID {
				temp = append(temp, node)
			}
		}
		AssignMatrixJoiningFee(temp, joining_package_fee, product_quantity)
		all_data = append(all_data, temp)
		cycleCount++

		for _, lst := range cycle {
			if len(lst) > 0 {
				allCycles = append(allCycles, lst)
			}
		}
	}
	return allCycles, cycleCount, all_data, totalJoiningFee
}

func CalculateMatrixSponsorBonus(allData [][]*TreeStructure, sponsorBonusPercent float64, mat_sponsor_bonus_type string, cappingAmount float64, cappingScope string) (float64, map[int]float64) {
	var TSB, sponsorBonus float64
	totalBonus := make(map[int]float64)
	for _, list := range allData {
		for _, member := range list {
			if member.ParentID != nil {
				if strings.Contains(mat_sponsor_bonus_type, "percent") {
					sponsorBonus = member.ParentID.SponsorBonus + member.Joining_package_fee*sponsorBonusPercent/100
				} else if strings.Contains(mat_sponsor_bonus_type, "usd") {
					sponsorBonus = member.ParentID.SponsorBonus + sponsorBonusPercent
				}
				if strings.Contains(cappingScope, "sponsor") && sponsorBonus > cappingAmount {
					member.ParentID.SponsorBonus = cappingAmount
				} else {
					member.ParentID.SponsorBonus = sponsorBonus
				}
				member.ParentID.SponsorBonus = sponsorBonus
			}
		}
	}

	for i, list := range allData {
		for _, member := range list {
			if member.UserID != 1 {
				TSB += member.SponsorBonus
				totalBonus[i] += member.SponsorBonus
			}
		}
	}
	return TSB, totalBonus
}

func CalculateMatrixMatchingBonus(allData [][]*TreeStructure, matchingPercentages []float64, cappingAmount float64, cappingScope string) (float64, map[int]float64) {
	var TMB float64
	totalBonus := make(map[int]float64)
	for _, members := range allData {
		for _, member := range members {
			iterant := 0
			if member.ParentID == nil {
				continue
			}
			parent := member.ParentID
			ApplyMatrixMatchingBonus(member, parent, matchingPercentages, iterant, cappingAmount, cappingScope)

		}
	}
	for i, list := range allData {
		for _, member := range list {
			if strings.Contains(cappingScope, "matching") && member.MatchingBonus > cappingAmount {
				member.MatchingBonus = cappingAmount
			}
			if member.UserID != 1 {
				TMB += member.MatchingBonus
				totalBonus[i] += member.MatchingBonus
			}
		}
	}
	return TMB, totalBonus
}

func ApplyMatrixMatchingBonus(member *TreeStructure, parent *TreeStructure, matchingPercentages []float64, iterant int, cappingAmount float64, cappingScope string) {
	if iterant >= len(matchingPercentages) || parent == nil {
		return
	}
	matching_bonus := parent.MatchingBonus + (member.SponsorBonus * matchingPercentages[iterant] / 100)

	parent.MatchingBonus = matching_bonus
	iterant = iterant + 1
	parent = parent.ParentID
	ApplyMatrixMatchingBonus(member, parent, matchingPercentages, iterant, cappingAmount, cappingScope)
}

func FindMatrixProfitToCompany(all_Data [][]*TreeStructure, ExpenseMembers float64, totalJoiningCycle map[int]float64, pool_perc float64) (float64, map[int]float64, float64, map[int]float64) {
	totalProfit := make(map[int]float64)
	var cycleProfits float64
	var pool_amount float64
	PoolAmount := make(map[int]float64)
	for i, list := range all_Data {
		var expense, cycleRevenue float64
		for _, member := range list {
			if member.UserID != 1 {
				cycleRevenue = cycleRevenue + member.Joining_package_fee
				expense = expense + ExpenseMembers + member.BinaryBonus + member.MatchingBonus + member.SponsorBonus + member.PoolBonus
			}
		}
		profit := cycleRevenue - expense
		expense = expense + (profit * (pool_perc / 100))

		PoolAmount[i] = profit * pool_perc / 100

		profit = cycleRevenue - expense
		cycleProfits = cycleProfits + profit
		totalProfit[i] = profit

		pool_amount = pool_amount + PoolAmount[i]
	}
	return cycleProfits, totalProfit, pool_amount, PoolAmount
}

func sendResultsToDjango(results interface{}) {
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:8000/process_results/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Response from Django:", resp.Status)
}

func main() {
	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// fmt.Println("Received data:", data)
		numMembers, ok := data["num_members"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'num_members' field", http.StatusBadRequest)
			return
		}
		PoolPerc, ok := data["pool_percentage"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'pool_percentage' field", http.StatusBadRequest)
			return
		}
		DistNo, ok := data["dist_member"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'dist_member' field", http.StatusBadRequest)
			return
		}
		ExpenseMembers, ok := data["expense_per_user"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'expense_per_user' field", http.StatusBadRequest)
			return
		}
		sponsorPercentage, ok := data["sponsor_percentage"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'sponsor_percentage' field", http.StatusBadRequest)
			return
		}
		binaryPercentage, ok := data["binary_percentage"].(float64)
		if !ok {
			binaryPercentage = -1.0
		}
		floatData := make([]float64, 0)
		percData := make([]float64, 0)
		bv := make([]float64, 0)

		var intData []int
		matchingPercentages, ok := data["matching_percentage"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'matching_percentage' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range matchingPercentages {
				num, ok := v.(int)
				if ok {
					percData = append(percData, float64(num))
					continue
				}
				numFloat64, ok := v.(float64)
				if ok {
					percData = append(percData, numFloat64)
					continue
				}
			}
		}
		joining_package_fee, ok := data["joining_package_fee"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'joining_package_fee' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range joining_package_fee {
				if num, ok := v.(float64); ok {
					floatData = append(floatData, float64(num))
				}
			}
		}
		b_volume, ok := data["b_volume"].([]interface{})
		if !ok {
			// http.Error(w, "Invalid or missing 'b_volume' field", http.StatusBadRequest)
			for range joining_package_fee {
				bv = append(bv, 0.0)
			}
			return
		} else {
			for _, v := range b_volume {
				if num, ok := v.(float64); ok {
					bv = append(bv, float64(num))
				}
			}
		}
		product_quantity, ok := data["product_quantity"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'product_quantity' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range product_quantity {
				if num, ok := v.(float64); ok {
					intData = append(intData, int(num))
				}
			}
		}
		cappingAmount, ok := data["capping_amount"].(float64)
		if !ok {
			cappingAmount = math.Pow(10, 100)
			return
		}
		cappingScope, ok := data["capping_scope"].(string)
		if !ok {
			http.Error(w, "Invalid or missing 'capping_scope' field", http.StatusBadRequest)
			return
		}
		bonusOption, ok := data["bonus_option"].(string)
		if !ok {
			http.Error(w, "Invalid or missing 'bonus_option' field", http.StatusBadRequest)
			return
		}
		sponsor_bonus_type1, ok := data["sponsor_bonus_type1"].(string)
		if !ok {
			http.Error(w, "Invalid or missing 'sponsor_bonus_type1' field", http.StatusBadRequest)
			return
		}

		sponsor_bonus_type2, ok := data["sponsor_bonus_type2"].(string)
		if !ok {
			http.Error(w, "Invalid or missing 'sponsor_bonus_type2' field", http.StatusBadRequest)
			return
		}
		binary_bonus_type, ok := data["binary_bonus_type"].(string)
		if !ok {
			http.Error(w, "Invalid or missing 'binary_bonus_type' field", http.StatusBadRequest)
			return
		}
		print(binary_bonus_type)
		binaryRatio, ok := data["ratio"].(string)
		if !ok {
			http.Error(w, "Missing or invalid 'ratio' field", http.StatusBadRequest)
			return
		}
		ratioAmount, ok := data["ratio_amount"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'ratio_amount' field", http.StatusBadRequest)
			return
		}

		root := &TreeStructure{UserID: 1, Levels: 0, Cycle: 1}
		result := buildTree(root, int(numMembers), 2)
		stored_id, cycleCount, all_data, totalJoiningCycle := AllocateMembers(int(numMembers), intData, 2, result, floatData, bv)
		var cycleList [][]*TreeStructure
		for _, list := range all_data {
			var copiedMembers []*TreeStructure
			for _, member := range list {
				copiedMember := &TreeStructure{
					UserID:              member.UserID,
					Levels:              member.Levels,
					ParentID:            member.ParentID,
					Cycle:               member.Cycle,
					Position:            member.Position,
					BV:                  member.BV,
					Joining_package_fee: member.Joining_package_fee,
					SponsorBonus:        member.SponsorBonus,
					BinaryBonus:         member.BinaryBonus,
					MatchingBonus:       member.MatchingBonus,
					CarryForward:        member.CarryForward,
					LeftMember:          member.LeftMember,
					RightMember:         member.RightMember,
					LeftCarry:           member.LeftCarry,
					RightCarry:          member.RightCarry,
					LeftDownlineSale:    member.LeftDownlineSale,
					RightDownlineSale:   member.RightDownlineSale,
					PoolBonus:           member.PoolBonus,
				}
				if copiedMember.LeftMember != nil {
					if copiedMember.LeftMember.UserID > len(list) {
						copiedMember.LeftMember = nil
					}
				}
				if copiedMember.RightMember != nil {
					if copiedMember.RightMember.UserID > len(list) {
						copiedMember.RightMember = nil
					}
				}
				copiedMembers = append(copiedMembers, copiedMember)
			}
			for _, member := range copiedMembers {
				if member.ParentID != nil {
					member.ParentID = copiedMembers[member.ParentID.UserID-1]
				}
			}
			cycleList = append(cycleList, copiedMembers)
		}

		totalSponsorBonus, totalSPONSORBonus := CalculateSponsorBonus(cycleList, sponsorPercentage, sponsor_bonus_type1, sponsor_bonus_type2, cappingAmount, cappingScope, bonusOption)
		totalBinaryBonus, totalBINARYBonus := BinaryWithRatio(cycleList, bonusOption, binaryRatio, int(ratioAmount), cappingScope, cappingAmount, cycleCount, binaryPercentage, binary_bonus_type)
		totalMatchingBonus, totalMATCHINGBonus := CalculateMatchingBonus(cycleList, percData, cappingAmount, cappingScope)
		totalProfitToCompany, Revenue, totalProfit, totalExpense, totalRevenue := FindProfitTOCompany(cycleList, ExpenseMembers, totalJoiningCycle)
		totalPoolBonus, totalPOOLBonus := CalculatePoolBonus(cycleList, totalProfit, PoolPerc, int(DistNo))
		TotalExpense := totalSponsorBonus + totalBinaryBonus + totalMatchingBonus + totalPoolBonus
		sum := 0.0
		for _, cycle := range totalProfit {
			sum = sum + cycle
		}
		totalProfitToCompany = sum
		results := map[string]interface{}{
			"totalProfitToCompany": totalProfitToCompany,
			"totalProfit":          totalProfit,
			"TotalExpense":         TotalExpense,
			"totalExpense":         totalExpense,
			"totalRevenue":         totalRevenue,
			"Revenue":              Revenue,
			// "treeNodes":            treeNodes,
			"stored_id":            stored_id,
			"cycleCount":           cycleCount,
			"total_sponsor_bonus":  totalSponsorBonus,
			"total_binary_bonus":   totalBinaryBonus,
			"total_matching_bonus": totalMatchingBonus,
			"totalSPONSORBonus":    totalSPONSORBonus,
			"totalBINARYBonus":     totalBINARYBonus,
			"totalMATCHINGBonus":   totalMATCHINGBonus,
			"totalPoolBonus":       totalPoolBonus,
			"totalPOOLBonus":       totalPOOLBonus,
		}
		sendResultsToDjango(results)

		fmt.Println("Results sent to Django.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	})
	http.HandleFunc("/unilevel", func(w http.ResponseWriter, r *http.Request) {
		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// fmt.Println("Received data:", data)
		numMembers, ok := data["num_members"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'num_members' field", http.StatusBadRequest)
			return
		}
		numChild, ok := data["num_child"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'num_child' field", http.StatusBadRequest)
			return
		}
		ExpenseMembers, ok := data["expense_per_user"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'expense_per_user' field", http.StatusBadRequest)
			return
		}
		sponsorPercentage, ok := data["sponsor_percentage"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'sponsor_percentage' field", http.StatusBadRequest)
			return
		}
		uni_sponsor_bonus_type, ok := data["uni_sponsor_bonus_type"].(string)
		if !ok {
			http.Error(w, "Invalid or missing 'uni_sponsor_bonus_type' field", http.StatusBadRequest)
			return
		}
		floatData := make([]float64, 0)
		percData := make([]float64, 0)

		var intData []int
		matchingPercentages, ok := data["matching_percentage"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'matching_percentage' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range matchingPercentages {
				num, ok := v.(int)
				if ok {
					percData = append(percData, float64(num))
					continue
				}
				numFloat64, ok := v.(float64)
				if ok {
					percData = append(percData, numFloat64)
					continue
				}
			}
		}
		joining_package_fee, ok := data["joining_package_fee"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'joining_package_fee' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range joining_package_fee {
				if num, ok := v.(float64); ok {
					floatData = append(floatData, float64(num))
				}
			}
		}
		product_quantity, ok := data["product_quantity"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'product_quantity' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range product_quantity {
				if num, ok := v.(float64); ok {
					intData = append(intData, int(num))
				}
			}
		}
		cappingAmount, ok := data["capping_amount"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'capping_amount' field", http.StatusBadRequest)
			return
		}
		cappingScope, ok := data["capping_scope"].(string)
		if !ok {
			http.Error(w, "Invalid or missing 'capping_scope' field", http.StatusBadRequest)
			return
		}
		PoolPerc, ok := data["pool_percentage"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'pool_percentage' field", http.StatusBadRequest)
			return
		}
		// DistNo, ok := data["dist_member"].(float64)
		// if !ok {
		// 	http.Error(w, "Invalid or missing 'dist_member' field", http.StatusBadRequest)
		// 	return
		// }
		root := &TreeStructure{UserID: 1, Levels: 0, Cycle: 1}
		result := buildUnilevelTree(root, int(numMembers), 2, int(numChild))
		stored_id, cycleCount, all_data, totalJoiningCycle := AllocateUnilevelMembers(int(numMembers), intData, 2, result, floatData)

		var cycleList [][]*TreeStructure
		for _, list := range all_data {
			var copiedMembers []*TreeStructure
			for _, member := range list {
				copiedMember := &TreeStructure{
					UserID:              member.UserID,
					Levels:              member.Levels,
					ParentID:            member.ParentID,
					Children:            member.Children,
					Cycle:               member.Cycle,
					Position:            member.Position,
					BV:                  member.BV,
					Joining_package_fee: member.Joining_package_fee,
					SponsorBonus:        member.SponsorBonus,
					BinaryBonus:         member.BinaryBonus,
					MatchingBonus:       member.MatchingBonus,
					CarryForward:        member.CarryForward,
					LeftMember:          member.LeftMember,
					RightMember:         member.RightMember,
					LeftCarry:           member.LeftCarry,
					RightCarry:          member.RightCarry,
					LeftDownlineSale:    member.LeftDownlineSale,
					RightDownlineSale:   member.RightDownlineSale,
					CycleProfit:         member.CycleProfit,
				}
				if copiedMember.LeftMember != nil {
					if copiedMember.LeftMember.UserID > len(list) {
						copiedMember.LeftMember = nil
					}
				}
				if copiedMember.RightMember != nil {
					if copiedMember.RightMember.UserID > len(list) {
						copiedMember.RightMember = nil
					}
				}
				if copiedMember.Children != nil {
					for _, child := range copiedMember.Children {
						if child.UserID > len(list) {
							child = nil
						}
					}
				}
				copiedMembers = append(copiedMembers, copiedMember)
			}
			for _, member := range copiedMembers {
				if member.ParentID != nil {
					member.ParentID = copiedMembers[member.ParentID.UserID-1]
				}
			}
			cycleList = append(cycleList, copiedMembers)
		}

		totalUnilevelSponsorBonus, totalUnilevelSPONSORBonus := CalculateUnilevelSponsorBonus(cycleList, sponsorPercentage, uni_sponsor_bonus_type, cappingAmount, cappingScope)
		totalUnilevelMatchingBonus, totalUnilevelMATCHINGBonus := CalculateUnilevelMatchingBonus(cycleList, percData, cappingAmount, cappingScope)
		totalUnilevelProfitToCompany, totalProfit, totalPoolBonus, totalPOOLBonus := FindUnilevelProfitToCompany(cycleList, ExpenseMembers, totalJoiningCycle, PoolPerc)
		// totalPoolBonus, totalPOOLBonus := CalculateUnilevelPoolBonus(cycleList, totalProfit, PoolPerc, int(DistNo))
		TotalUnilevelExpense := totalUnilevelSponsorBonus + totalUnilevelMatchingBonus + totalPoolBonus
		// var treeNodes [][]TreeStructureJSON
		// for _, list := range cycleList {
		// 	temp := convertToJSONStructure(list)
		// 	treeNodes = append(treeNodes, temp)
		// }

		results := map[string]interface{}{
			"totalPoolBonus":       totalPoolBonus,
			"totalPOOLBonus":       totalPOOLBonus,
			"TotalUnilevelExpense": TotalUnilevelExpense,
			"totalProfit":          totalProfit,
			"totalJoiningCycle":    totalJoiningCycle,
			// "treeNodes":                    treeNodes,
			"stored_id":                    stored_id,
			"cycleCount":                   cycleCount,
			"totalUnilevelProfitToCompany": totalUnilevelProfitToCompany,
			"totalUnilevelSponsorBonus":    totalUnilevelSponsorBonus,
			"totalUnilevelMatchingBonus":   totalUnilevelMatchingBonus,
			"totalUnilevelSPONSORBonus":    totalUnilevelSPONSORBonus,
			"totalUnilevelMATCHINGBonus":   totalUnilevelMATCHINGBonus,
		}
		sendResultsToDjango(results)

		fmt.Println("Results sent to Django.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	})

	http.HandleFunc("/matrix", func(w http.ResponseWriter, r *http.Request) {
		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// fmt.Println("Received data:", data)
		numMembers, ok := data["num_members"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'num_members' field", http.StatusBadRequest)
			return
		}
		numChild, ok := data["num_child"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'num_child' field", http.StatusBadRequest)
			return
		}
		ExpenseMembers, ok := data["expense_per_user"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'expense_per_user' field", http.StatusBadRequest)
			return
		}
		fmt.Println("11")
		sponsorPercentage, ok := data["sponsor_percentage"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'sponsor_percentage' field", http.StatusBadRequest)
			return
		}
		fmt.Println("22")
		mat_sponsor_bonus_type, ok := data["mat_sponsor_bonus_type"].(string)
		if !ok {
			http.Error(w, "Invalid or missing 'mat_sponsor_bonus_type' field", http.StatusBadRequest)
			return
		}
		floatData := make([]float64, 0)
		percData := make([]float64, 0)

		var intData []int
		matchingPercentages, ok := data["matching_percentage"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'matching_percentage' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range matchingPercentages {
				num, ok := v.(int)
				if ok {
					percData = append(percData, float64(num))
					continue
				}
				numFloat64, ok := v.(float64)
				if ok {
					percData = append(percData, numFloat64)
					continue
				}
			}
		}
		PoolPerc, ok := data["pool_percentage"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'pool_percentage' field", http.StatusBadRequest)
			return
		}
		joining_package_fee, ok := data["joining_package_fee"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'joining_package_fee' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range joining_package_fee {
				if num, ok := v.(float64); ok {
					floatData = append(floatData, float64(num))
				}
			}
		}
		product_quantity, ok := data["product_quantity"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'product_quantity' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range product_quantity {
				if num, ok := v.(float64); ok {
					intData = append(intData, int(num))
				}
			}
		}
		cappingAmount, ok := data["capping_amount"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'capping_amount' field", http.StatusBadRequest)
			return
		}
		cappingScope, ok := data["capping_scope"].(string)
		if !ok {
			http.Error(w, "Invalid or missing 'capping_scope' field", http.StatusBadRequest)
			return
		}
		root := &TreeStructure{UserID: 1, Levels: 0, Cycle: 1}
		result := buildUnilevelTree(root, int(numMembers), 2, int(numChild))
		stored_id, cycleCount, all_data, totalJoiningCycle := AllocateMatrixMembers(int(numMembers), intData, 2, result, floatData)

		var cycleList [][]*TreeStructure
		for _, list := range all_data {
			var copiedMembers []*TreeStructure
			for _, member := range list {
				copiedMember := &TreeStructure{
					UserID:              member.UserID,
					Levels:              member.Levels,
					ParentID:            member.ParentID,
					Children:            member.Children,
					Cycle:               member.Cycle,
					Position:            member.Position,
					BV:                  member.BV,
					Joining_package_fee: member.Joining_package_fee,
					SponsorBonus:        member.SponsorBonus,
					BinaryBonus:         member.BinaryBonus,
					MatchingBonus:       member.MatchingBonus,
					CarryForward:        member.CarryForward,
					LeftMember:          member.LeftMember,
					RightMember:         member.RightMember,
					LeftCarry:           member.LeftCarry,
					RightCarry:          member.RightCarry,
					LeftDownlineSale:    member.LeftDownlineSale,
					RightDownlineSale:   member.RightDownlineSale,
				}
				if copiedMember.LeftMember != nil {
					if copiedMember.LeftMember.UserID > len(list) {
						copiedMember.LeftMember = nil
					}
				}
				if copiedMember.RightMember != nil {
					if copiedMember.RightMember.UserID > len(list) {
						copiedMember.RightMember = nil
					}
				}
				if copiedMember.Children != nil {
					for _, child := range copiedMember.Children {
						if child.UserID > len(list) {
							child = nil
						}
					}
				}
				copiedMembers = append(copiedMembers, copiedMember)
			}
			for _, member := range copiedMembers {
				if member.ParentID != nil {
					member.ParentID = copiedMembers[member.ParentID.UserID-1]
				}
			}
			cycleList = append(cycleList, copiedMembers)
		}

		CalculateMatrixSponsorBonus, totalMatrixSPONSORBonus := CalculateMatrixSponsorBonus(cycleList, sponsorPercentage, mat_sponsor_bonus_type, cappingAmount, cappingScope)
		CalculateMatrixMatchingBonus, totalMatrixMATCHINGBonus := CalculateMatrixMatchingBonus(cycleList, percData, cappingAmount, cappingScope)
		FindMatrixProfitToCompany, totalProfit, totalPoolBonus, totalPOOLBonus := FindMatrixProfitToCompany(cycleList, ExpenseMembers, totalJoiningCycle, PoolPerc)
		// totalPoolBonus, totalPOOLBonus := CalculateMatrixPoolBonus(cycleList, totalProfit, PoolPerc, int(DistNo))
		TotalMatrixExpense := CalculateMatrixSponsorBonus + CalculateMatrixMatchingBonus + totalPoolBonus

		results := map[string]interface{}{
			"totalPoolBonus":     totalPoolBonus,
			"totalPOOLBonus":     totalPOOLBonus,
			"TotalMatrixExpense": TotalMatrixExpense,
			"totalJoiningCycle":  totalJoiningCycle,
			"totalProfit":        totalProfit,
			// "treeNodes":                    treeNodes,
			"stored_id":                    stored_id,
			"cycleCount":                   cycleCount,
			"totalMatrixProfitToCompany":   FindMatrixProfitToCompany,
			"CalculateMatrixSponsorBonus":  CalculateMatrixSponsorBonus,
			"CalculateMatrixMatchingBonus": CalculateMatrixMatchingBonus,
			"totalMatrixSPONSORBonus":      totalMatrixSPONSORBonus,
			"totalMatrixMATCHINGBonus":     totalMatrixMATCHINGBonus,
		}
		sendResultsToDjango(results)

		// fmt.Println("Results sent to Django.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}
