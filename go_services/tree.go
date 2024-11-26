package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

type TreeStructure struct {
	UserID              int
	ParentID            *TreeStructure
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
	FlushOut            int
	LeftCarry           int
	RightCarry          int
	LeftDownlineSale    float64
	RightDownlineSale   float64
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
	FlushOut            int     `json:"flush_out"`
	LeftCarry           int     `json:"left_carry"`
	RightCarry          int     `json:"right_carry"`
	LeftDownlineSale    float64 `json:"LeftDownlineSale"`
	RightDownlineSale   float64 `json:"RightDownlineSale"`
}

func convertToJSONStructure(nodes []*TreeStructure) []TreeStructureJSON {
	var jsonNodes []TreeStructureJSON
	for _, node := range nodes {
		var parentID int
		var leftMemberID, rightMemberID int
		if node.ParentID != nil {
			parentID = node.ParentID.UserID
		} else {
			parentID = 0
		}
		if node.LeftMember != nil {
			leftMemberID = node.LeftMember.UserID
		} else {
			leftMemberID = 0
		}

		if node.RightMember != nil {
			rightMemberID = node.RightMember.UserID
		} else {
			rightMemberID = 0
		}

		jsonNodes = append(jsonNodes, TreeStructureJSON{
			UserID:              node.UserID,
			Position:            node.Position,
			Joining_package_fee: node.Joining_package_fee,
			BV:                  node.BV,
			LeftMember:          leftMemberID,
			RightMember:         rightMemberID,
			Levels:              node.Levels,
			ParentID:            parentID,
			SponsorBonus:        node.SponsorBonus,
			BinaryBonus:         node.BinaryBonus,
			MatchingBonus:       node.MatchingBonus,
			CarryForward:        node.CarryForward,
			Cycle:               node.Cycle,
			FlushOut:            node.FlushOut,
			LeftCarry:           node.LeftCarry,
			RightCarry:          node.RightCarry,
			LeftDownlineSale:    node.LeftDownlineSale,
			RightDownlineSale:   node.RightDownlineSale,
		})
	}
	return jsonNodes
}

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

func AllocateMembers(numMembers int, product_quantity []int, startID int, result []*TreeStructure, joining_package_fee []float64, b_volume []float64) ([][]int, int, [][]*TreeStructure) {
	currentID := startID
	remaing := numMembers + 2
	var allCycles [][]int
	cycleCount := 1
	for remaing > 2 {
		var cycle [][]int
		for _, qty := range product_quantity {
			cycleID := []int{}
			for j := 0; j < qty && remaing > 2; j++ {
				cycleID = append(cycleID, currentID)
				currentID++
				remaing--
			}
			cycle = append(cycle, cycleID)
			if remaing <= 0 {
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
	return allCycles, cycleCount, all_data
}

func CalculateSponsorBonus(allData [][]*TreeStructure, sponsorBonusPercent float64, cappingAmount float64, cappingScope, bonusOption string) map[int]float64 {
	totalBonus := make(map[int]float64)

	calculateBonus := func(member *TreeStructure, valueSelector func(*TreeStructure) float64) float64 {
		var rightBonus, leftBonus float64

		if member.RightMember != nil {
			rightBonus = valueSelector(member.RightMember) * sponsorBonusPercent / 100
		}
		if member.LeftMember != nil {
			leftBonus = valueSelector(member.LeftMember) * sponsorBonusPercent / 100
		}

		totalBonus := rightBonus + leftBonus
		if strings.Contains(cappingScope, "sponsor") && totalBonus > cappingAmount {
			print(totalBonus)
			return cappingAmount
		}
		return totalBonus
	}

	for i, list := range allData {
		for _, member := range list {
			var bonus float64

			if strings.Contains(bonusOption, "PRICE") {
				bonus = calculateBonus(member, func(m *TreeStructure) float64 { return m.Joining_package_fee })
			} else if strings.Contains(bonusOption, "BV") {
				bonus = calculateBonus(member, func(m *TreeStructure) float64 { return m.BV })
			}

			member.SponsorBonus = bonus
			totalBonus[i] += bonus
		}
	}

	return totalBonus
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

func BinaryWithRatio(allData [][]*TreeStructure, bonusOption string, binaryRatio string, ratioAmount int, cappingScope string, cappingAmount float64, cycleCount int) map[int]float64 {
	start := time.Now()
	total := make(map[int]float64)
	b1 := 10
	b2 := 15
	b3 := 20
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

			var left, right, noOfPairs, carryLeft, carryRight int
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
			carryLeft = int(node.LeftDownlineSale) - noOfPairs*left
			carryRight = int(node.RightDownlineSale) - noOfPairs*right

			node.LeftCarry = carryLeft
			node.RightCarry = carryRight

			minimumVal := int(math.Min(float64(noOfPairs*left), float64(noOfPairs*right)))

			var bonusPerc int
			switch {
			case noOfPairs > 0 && noOfPairs <= 5:
				bonusPerc = b1
			case noOfPairs > 5 && noOfPairs <= 10:
				bonusPerc = b2
			case noOfPairs > 10:
				bonusPerc = b3
			default:
				bonusPerc = 0
			}

			nodeBonus := (minimumVal * bonusPerc) / 100

			if strings.Contains(cappingScope, "binary") {
				if nodeBonus > int(cappingAmount) {
					node.FlushOut = nodeBonus - int(cappingAmount)
					nodeBonus = int(cappingAmount)
				}
			}

			node.BinaryBonus = float64(nodeBonus)
			total[i] += node.BinaryBonus

		}
	}
	elapsed := time.Since(start)
	fmt.Println("time elapse in nanoseconds:", elapsed)
	return total
}

func CalculateMatchingBonus(allData [][]*TreeStructure, matchingPercentages []float64, cappingAmount float64, cappingScope string) map[int]float64 {
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
			totalBonus[i] += member.MatchingBonus
		}
	}
	return totalBonus
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

		fmt.Println("Received data:", data)

		numMembers, ok := data["num_members"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'num_members' field", http.StatusBadRequest)
			return
		}
		sponsorPercentage, ok := data["sponsor_percentage"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'sponsor_percentage' field", http.StatusBadRequest)
			return
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
			http.Error(w, "Invalid or missing 'b_volume' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range b_volume {
				if num, ok := v.(float64); ok {
					bv = append(bv, float64(num))
				}
			}
		}
		// productName, ok := data["product_name"].([]interface{})
		// if !ok {
		// 	http.Error(w, "Invalid or missing 'product_name' field", http.StatusBadRequest)
		// 	return
		// }
		// var p_name []string
		// for _, v := range productName {
		// 	if str, ok := v.(string); ok {
		// 		p_name = append(p_name, str)
		// 	}
		// }
		// print("---------------", productName)
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
		//scopes := strings.Split(cappingScope, ",")
		fmt.Println(".............................", cappingScope)

		bonusOption, ok := data["bonus_option"].(string)
		if !ok {
			http.Error(w, "Invalid or missing 'bonus_option' field", http.StatusBadRequest)
			return
		}
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

		stored_id, cycleCount, all_data := AllocateMembers(int(numMembers), intData, 2, result, floatData, bv)
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

		totalSponsorBonus := CalculateSponsorBonus(cycleList, sponsorPercentage, cappingAmount, cappingScope, bonusOption)
		totalBinaryBonus := BinaryWithRatio(cycleList, bonusOption, binaryRatio, int(ratioAmount), cappingScope, cappingAmount, cycleCount)
		totalMatchingBonus := CalculateMatchingBonus(cycleList, percData, cappingAmount, cappingScope)
		var treeNodes [][]TreeStructureJSON
		for _, list := range cycleList {
			temp := convertToJSONStructure(list)
			treeNodes = append(treeNodes, temp)
		}
		fmt.Println("000", stored_id)
		results := map[string]interface{}{
			"treeNodes":            treeNodes,
			"stored_id":            stored_id,
			"cycleCount":           cycleCount,
			"total_sponsor_bonus":  totalSponsorBonus,
			"total_binary_bonus":   totalBinaryBonus,
			"total_matching_bonus": totalMatchingBonus,
		}
		sendResultsToDjango(results)

		fmt.Println("Results sent to Django.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	})
	log.Fatal(http.ListenAndServe(":9000", nil))
}

// func CalculateSponsorBonus(allData [][]*TreeStructure, sponsorBonusPercent float64, joining_package_fee []float64, cappingAmount float64, cappingScope string, numMembers int, bonusOption string) map[int]float64 {
// 	totalBonus := make(map[int]float64)
// 	currentCycle := 1
// 	for i, list := range allData {
// 		for _, member := range list {
// 			var rightBonus, leftBonus float64
// 			if strings.Contains(bonusOption, "PRICE"){
// 				if member.RightMember != nil {
// 					rightBonus = member.RightMember.Joining_package_fee * sponsorBonusPercent / 100
// 				}
// 				if member.LeftMember != nil {
// 					leftBonus = member.LeftMember.Joining_package_fee * sponsorBonusPercent / 100
// 				}
// 				sponsorBonus := rightBonus + leftBonus
// 				if strings.Contains(cappingScope, "sponsor") && sponsorBonus > cappingAmount {
// 					member.SponsorBonus = cappingAmount
// 				} else {
// 					member.SponsorBonus = sponsorBonus
// 				}
// 			}
// 			if strings.Contains(bonusOption, "BV"){
// 				if member.RightMember != nil {
// 					rightBonus = member.RightMember.BV * sponsorBonusPercent / 100
// 				}
// 				if member.LeftMember != nil {
// 					leftBonus = member.LeftMember.BV * sponsorBonusPercent / 100
// 				}
// 				sponsorBonus := rightBonus + leftBonus
// 				if strings.Contains(cappingScope, "sponsor") && sponsorBonus > cappingAmount {
// 					member.SponsorBonus = cappingAmount
// 				} else {
// 					member.SponsorBonus = sponsorBonus
// 				}
// 			}
// 			totalBonus[i] += member.SponsorBonus

// 		}
// 		currentCycle = currentCycle + 1
// 	}
// 	return totalBonus
// }
