package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	// "strconv"
	"math"
	"net/http"
	"strings"
)

type TreeStructure struct {
	UserID              int
	ParentID            *TreeStructure
	Children            []*TreeStructure
	Joining_package_fee float64
	Position            string
	Levels              int
	Lft                 int
	Rgt                 int
	LeftMember          *TreeStructure
	RightMember         *TreeStructure
	SponsorBonus        float64
	BinaryBonus         float64
	LeftSales           float64
	RightSales          float64
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
	UserID              int                 `json:"user_id"`
	Position            string              `json:"position"`
	Joining_package_fee float64             `json:"joining_package_fee"`
	Levels              int                 `json:"levels"`
	Lft                 int                 `json:"lft"`
	Rgt                 int                 `json:"rgt"`
	ParentID            int                 `json:"parentid"`
	SponsorBonus        float64             `json:"sponsor_bonus"`
	BinaryBonus         float64             `json:"binary_bonus"`
	MatchingBonus       float64             `json:"matching_bonus"`
	CarryForward        float64             `json:"carry_forward"`
	Cycle               int                 `json:"cycle"`
	FlushOut            int                 `json:"flush_out"`
	LeftCarry           int                 `json:"left_carry"`
	RightCarry          int                 `json:"right_carry"`
	Children            []TreeStructureJSON `json:"children"`
}

func convertToJSONStructure(nodes []*TreeStructure) []TreeStructureJSON {
	var jsonNodes []TreeStructureJSON
	for _, node := range nodes {
		var parentID int
		if node.ParentID != nil {
			parentID = node.ParentID.UserID
		} else {
			parentID = 0
		}

		jsonNodes = append(jsonNodes, TreeStructureJSON{
			UserID:              node.UserID,
			Position:            node.Position,
			Joining_package_fee: node.Joining_package_fee,
			Levels:              node.Levels,
			ParentID:            parentID,
			SponsorBonus:        node.SponsorBonus,
			BinaryBonus:         node.BinaryBonus,
			MatchingBonus:       node.MatchingBonus,
			CarryForward:        node.CarryForward,
			Cycle:               node.Cycle,
		})
	}
	return jsonNodes
}

var treeNodes []*TreeStructure

func buildTree(root *TreeStructure, numMembers int, startID int, joining_package_fee []float64) []*TreeStructure {
	if numMembers <= 0 {
		fmt.Println("Invalid input: numMembers or cycles is non-positive")
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

func AllocateMembers(numMembers int, product_quantity []int, startID int, result []*TreeStructure, joining_package_fee []float64) ([][]int, int, [][]*TreeStructure) {
	currentID := startID
	remaing := numMembers + 2
	var allCycles [][]int
	cycleCount := 0
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
		AssignJoiningFee(temp, joining_package_fee, product_quantity)
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

func AssignJoiningFee(nodes []*TreeStructure, joining_package_fee []float64, product_quantity []int) {
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

func CalculateSponsorBonus(nodes []*TreeStructure, sponsorBonusPercent float64, joining_package_fee []float64, cappingAmount float64, cappingScope string, numMembers int) float64 {
	var totalBonus float64
	currentCycle := 1
	var cycleMemberCount int

	for _, member := range nodes {
		var rightBonus, leftBonus float64

		if member.RightMember != nil {
			rightBonus = member.RightMember.Joining_package_fee * sponsorBonusPercent / 100
		}
		if member.LeftMember != nil {
			leftBonus = member.LeftMember.Joining_package_fee * sponsorBonusPercent / 100
		}

		sponsorBonus := rightBonus + leftBonus
		if strings.Contains(cappingScope, "sponsor") && sponsorBonus > cappingAmount {
			member.SponsorBonus = cappingAmount
		} else {
			member.SponsorBonus = sponsorBonus
		}
		totalBonus += member.SponsorBonus
		cycleMemberCount++
		if cycleMemberCount == numMembers {
			cycleMemberCount = 0
			currentCycle++
		}
	}
	return totalBonus
}

func Traverse(node *TreeStructure, limit int) float64 {
	if node == nil || node.UserID > limit {
		return 0
	}
	currentSales := node.Joining_package_fee
	leftSales := Traverse(node.LeftMember, limit)
	rightSales := Traverse(node.RightMember, limit)
	return currentSales + leftSales + rightSales
}

func BinaryWithRatio(allData [][]*TreeStructure, joining_package_fee []float64, binaryRatio string, ratioAmount int, cappingScope string, cappingAmount float64) float64 {

	var total float64
	b1 := 10
	b2 := 15
	b3 := 20
	fmt.Println("Starting BinaryWithRatio function...")
	for _, nodeList := range allData {
		for _, node := range nodeList {
			var leftDownline, rightDownline float64
			if node.LeftMember != nil {
				leftDownline = Traverse(node.LeftMember, len(nodeList))
			}
			if node.RightMember != nil {
				rightDownline = Traverse(node.RightMember, len(nodeList))
			}
			node.LeftDownlineSale = float64(leftDownline)
			node.RightDownlineSale = float64(rightDownline)

			fmt.Printf("Node %d: LeftDownlineSale: %f, RightDownlineSale: %f\n", node.UserID, node.LeftDownlineSale, node.RightDownlineSale)

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
			if cappingScope == "binary" {
				if nodeBonus > int(cappingAmount) {
					node.FlushOut = nodeBonus - int(cappingAmount)
					nodeBonus = int(cappingAmount)
					fmt.Printf("Node %d: Bonus capped at %f, FlushOut: %d\n", node.UserID, cappingAmount, node.FlushOut)
				}
			}
			node.BinaryBonus = float64(nodeBonus)
			total += float64(nodeBonus)
		}
	}

	return total
}

func CalculateMatchingBonus(members []*TreeStructure, matchingPercentages []float64, cappingAmount float64, cappingScope string) float64 {
	var totalBonus float64
	for _, member := range members {
		if member.ParentID == nil {
			continue
		}
		parent := member.ParentID

		if parent.MatchingBonus == 0 {
			parent.MatchingBonus = 0
		}

		ApplyMatchingBonus(member, parent, matchingPercentages, 0, cappingAmount, cappingScope)
	}

	for _, member := range members {
		totalBonus += member.MatchingBonus
	}

	return totalBonus
}

func ApplyMatchingBonus(member, parent *TreeStructure, matchingPercentages []float64, iterant int, cappingAmount float64, cappingScope string) {

	if iterant >= len(matchingPercentages) || parent == nil {

		return
	}

	if member.BinaryBonus == 0 {
		member.BinaryBonus = 0
	}

	matchingBonus := parent.MatchingBonus + (member.BinaryBonus * matchingPercentages[iterant] / 100)

	if strings.Contains(cappingScope, "matching") && matchingBonus > cappingAmount {
		parent.MatchingBonus = cappingAmount
	} else {
		parent.MatchingBonus = matchingBonus
	}

	ApplyMatchingBonus(member, parent.ParentID, matchingPercentages, iterant+1, cappingAmount, cappingScope)
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
		// binaryPercentage, ok := data["binary_percentage"].(float64)
		// if !ok {
		// 	http.Error(w, "Invalid or missing 'binary_percentage' field", http.StatusBadRequest)
		// 	return
		// }
		floatData := make([]float64, 0)
		var intData []int
		matchingPercentages, ok := data["joining_package_fee"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'joining_package_fee' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range matchingPercentages {
				num, ok := v.(int)
				if ok {
					floatData = append(floatData, float64(num))
					continue
				}
				numFloat64, ok := v.(float64)
				if ok {
					floatData = append(floatData, numFloat64)
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
				if num, ok := v.(int); ok {
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

		var treeNodes []*TreeStructure
		root := &TreeStructure{UserID: 1, Levels: 0, Cycle: 1}
		result := buildTree(root, int(numMembers), 2, floatData)
		for _, tree := range result {
			treeNodes = append(treeNodes, tree)
		}
		stored_id, cycleCount, all_data := AllocateMembers(int(numMembers), intData, 2, result, floatData)
		//AssignJoiningFee(treeNodes, stored_id, floatData)
		// for _, node := range treeNodes {
		// 	fmt.Println(node)
		// }
		for _, list := range all_data {
			for _, data := range list {
				fmt.Println("^^^^^^^^^^^^^^^^", data.Joining_package_fee)
			}
		}

		totalSponsorBonus := CalculateSponsorBonus(treeNodes, sponsorPercentage, floatData, cappingAmount, cappingScope, int(numMembers))
		cycle := 1
		for _, list := range all_data {
			fmt.Println(cycle)
			for _, member := range list {
				fmt.Println("222222222", member.Joining_package_fee)
			}
			cycle++
		}
		totalBinaryBonus := BinaryWithRatio(all_data, floatData, binaryRatio, int(ratioAmount), cappingScope, cappingAmount)
		// totalBinaryBonus := binaryWithRatio(all_data, binaryRatio, ratioAmount, cappingScope, cappingAmount)

		totalMatchingBonus := CalculateMatchingBonus(treeNodes, floatData, cappingAmount, cappingScope)
		var temp []TreeStructureJSON
		var all_data_json [][]TreeStructureJSON
		for _, list := range all_data {
			temp = convertToJSONStructure(list)
			all_data_json = append(all_data_json, temp)
		}
		results := map[string]interface{}{
			"tree_structure":       convertToJSONStructure(treeNodes),
			"stored_id":            stored_id,
			"cycleCount":           cycleCount,
			"all_data":             all_data_json,
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

/*
// func binaryWithRatio(allData [][]*TreeStructure, binaryRatio string, ratioAmount float64, cappingChoice string, cappingValue float64) float64 {
// 	var totalBonus float64
// 	b1 := 5
// 	b2 := 10
// 	b3 := 15
// 	fmt.Println("Starting calculation with binaryRatio:", binaryRatio, "ratioAmount:", ratioAmount, "cappingChoice:", cappingChoice, "cappingValue:", cappingValue)

// 	for _, list := range allData {
// 		for _, i := range list {
// 			var left, right float64
// 			switch binaryRatio {
// 			case "1":
// 				left = ratioAmount * 1
// 				right = ratioAmount * 1
// 			case "2":
// 				left = ratioAmount * 1
// 				right = ratioAmount * 2
// 			case "3":
// 				left = ratioAmount * 2
// 				right = ratioAmount * 1
// 			}
// 			fmt.Println("Left:", left, "Right:", right)

// 			left_sale := 0.0
// 			right_sale := 0.0
// 			if i.LeftMember != nil {
// 				left_sale = Traverse(i.LeftMember)
// 				fmt.Println("Left sale:", left_sale)
// 			}
// 			if i.RightMember != nil {
// 				right_sale = Traverse(i.RightMember)
// 				fmt.Println("Right sale:", right_sale)
// 			}

// 			leftTemp := int(math.Floor(float64(left_sale) / left))
// 			rightTemp := int(math.Floor(float64(right_sale) / right))
// 			noOfPairs := int(math.Min(float64(leftTemp), float64(rightTemp)))
// 			fmt.Println("Left Temp:", leftTemp, "Right Temp:", rightTemp, "No of Pairs:", noOfPairs)

// 			var bonusPercentage float64
// 			if noOfPairs <= 5 {
// 				bonusPercentage = float64(b1)
// 			} else if noOfPairs >= 6 && noOfPairs <= 10 {
// 				bonusPercentage = float64(b2)
// 			} else {
// 				bonusPercentage = float64(b3)
// 			}
// 			fmt.Println("Bonus Percentage:", bonusPercentage)

// 			minVal := math.Min(float64(noOfPairs)*left, float64(noOfPairs)*right)
// 			totalBonusForUser := (minVal * bonusPercentage) / 100
// 			i.BinaryBonus = totalBonusForUser
// 			fmt.Println("Total bonus for user:", totalBonusForUser)

// 			if cappingChoice == "binary" {
// 				totalBonusForUser = math.Min(totalBonusForUser, cappingValue)
// 				fmt.Println("Capped total bonus for user:", totalBonusForUser)
// 			}
// 			totalBonus += totalBonusForUser
// 		}
// 	}

// 	fmt.Println("Total bonus calculated:", totalBonus)
// 	return totalBonus
// }


type TreeStructure struct {
	UserID              int
	ParentID            *TreeStructure
	Children            []*TreeStructure
	LeftCarry           int
	RightCarry          int
	Joining_package_fee float64
	Position            string
	Levels              int
	Lft                 int
	Rgt                 int
	LeftMember          *TreeStructure
	RightMember         *TreeStructure
	SponsorBonus        float64
	BinaryBonus         float64
	LeftSales           float64
	RightSales          float64
	CarryForward        float64
	MatchingBonus       float64
	CappingValue        float64
	Cycle               int
	LeftDownlineSale    float64
	RightDownlineSale   float64
	FlushOut            int
}

type TreeStructureJSON struct {
	UserID              int                 `json:"user_id"`
	Position            string              `json:"position"`
	Joining_package_fee float64             `json:"joining_package_fee"`
	Levels              int                 `json:"levels"`
	Lft                 int                 `json:"lft"`
	Rgt                 int                 `json:"rgt"`
	ParentID            int                 `json:"parentid"`
	SponsorBonus        float64             `json:"sponsor_bonus"`
	BinaryBonus         float64             `json:"binary_bonus"`
	MatchingBonus       float64             `json:"matching_bonus"`
	CarryForward        float64             `json:"carry_forward"`
	Cycle               int                 `json:"cycle"`
	Children            []TreeStructureJSON `json:"children"`
}

func convertToJSONStructure(nodes []*TreeStructure) []TreeStructureJSON {
	var jsonNodes []TreeStructureJSON
	for _, node := range nodes {
		var parentID int
		if node.ParentID != nil {
			parentID = node.ParentID.UserID
		} else {
			parentID = 0
		}
		var children []TreeStructureJSON
		if len(node.Children) > 0 {
			children = convertToJSONStructure(node.Children)
		}
		jsonNodes = append(jsonNodes, TreeStructureJSON{
			UserID:              node.UserID,
			Position:            node.Position,
			Joining_package_fee: node.Joining_package_fee,
			Levels:              node.Levels,
			ParentID:            parentID,
			SponsorBonus:        node.SponsorBonus,
			BinaryBonus:         node.BinaryBonus,
			MatchingBonus:       node.MatchingBonus,
			CarryForward:        node.CarryForward,
			Cycle:               node.Cycle,
			Children:            children,
		})
	}
	return jsonNodes
}

var treeNodes []*TreeStructure





func Traverse(node *TreeStructure) float64 {
	if node == nil {
		return 0
	}
	currentSales := node.Joining_package_fee

	leftSales := Traverse(node.LeftMember)
	rightSales := Traverse(node.RightMember)

	return currentSales + leftSales + rightSales
}

// func BinaryWithRatio(allData [][]*TreeStructure, joining_package_fee []float64, binaryRatio string, ratioAmount int, cappingScope string, cappingAmount float64) float64 {

// 	var total float64
// 	b1 := 5
// 	b2 := 10
// 	b3 := 15

// 	fmt.Println("Starting BinaryWithRatio function...")
// 	for _, nodeList := range allData {

// 		fmt.Println(len(nodeList))
// 		for _, node := range nodeList {
// 			if len(node.Children) < 2 {
// 				fmt.Println("Skipping node (not enough children):", node.UserID)
// 				continue
// 			}
// 			leftChild := node.Children[0]
// 			rightChild := node.Children[1]

// 			leftDownline := Traverse(leftChild)
// 			rightDownline := Traverse(rightChild)

// 			node.LeftDownlineSale = float64(leftDownline) + float64(node.LeftCarry)
// 			node.RightDownlineSale = float64(rightDownline) + float64(node.RightCarry)

// 			fmt.Printf("Node %d: LeftDownlineSale: %f, RightDownlineSale: %f\n", node.UserID, node.LeftDownlineSale, node.RightDownlineSale)

// 			var left, right, noOfPairs, carryLeft, carryRight int
// 			switch binaryRatio {
// 			case "1":
// 				left, right = ratioAmount*1, ratioAmount*1
// 			case "2":
// 				left, right = ratioAmount*1, ratioAmount*2
// 			default:
// 				left, right = ratioAmount*2, ratioAmount*1
// 			}

// 			fmt.Printf("Node %d: Left Ratio: %d, Right Ratio: %d\n", node.UserID, left, right)

// 			leftTemp := int(math.Floor(float64(node.LeftDownlineSale) / float64(left)))
// 			rightTemp := int(math.Floor(float64(node.RightDownlineSale) / float64(right)))

// 			noOfPairs = int(math.Min(float64(leftTemp), float64(rightTemp)))
// 			carryLeft = int(node.LeftDownlineSale) - noOfPairs*left
// 			carryRight = int(node.RightDownlineSale) - noOfPairs*right

// 			node.LeftCarry = carryLeft
// 			node.RightCarry = carryRight

// 			fmt.Printf("Node %d: NoOfPairs: %d, LeftCarry: %d, RightCarry: %d\n", node.UserID, noOfPairs, carryLeft, carryRight)

// 			minimumVal := int(math.Min(float64(noOfPairs*left), float64(noOfPairs*right)))

// 			var bonusPerc int
// 			switch {
// 			case noOfPairs > 0 && noOfPairs <= 5:
// 				bonusPerc = b1
// 			case noOfPairs > 5 && noOfPairs <= 10:
// 				bonusPerc = b2
// 			case noOfPairs > 10:
// 				bonusPerc = b3
// 			default:
// 				bonusPerc = 0
// 			}

// 			fmt.Printf("Node %d: Bonus Percentage: %d\n", node.UserID, bonusPerc)

// 			nodeBonus := (minimumVal * bonusPerc) / 100
// 			if cappingScope == "binary" {
// 				if nodeBonus > int(cappingAmount) {
// 					node.FlushOut = nodeBonus - int(cappingAmount)
// 					nodeBonus = int(cappingAmount)
// 					fmt.Printf("Node %d: Bonus capped at %f, FlushOut: %d\n", node.UserID, cappingAmount, node.FlushOut)
// 				}
// 			}

// 			node.BinaryBonus = float64(nodeBonus)
// 			total += float64(nodeBonus)

// 			fmt.Printf("Node %d: BinaryBonus: %f, Total: %f\n", node.UserID, node.BinaryBonus, total)
// 		}
// 	}

// 	fmt.Println("BinaryWithRatio function complete.")
// 	return total
// }




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
		binaryPercentage, ok := data["binary_percentage"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing 'binary_percentage' field", http.StatusBadRequest)
			return
		}
		floatData := make([]float64, 0)
		var intData []int
		matchingPercentages, ok := data["matching_percentage"].([]interface{})
		if !ok {
			http.Error(w, "Invalid or missing 'matching_percentage' field", http.StatusBadRequest)
			return
		} else {
			for _, v := range matchingPercentages {
				num, ok := v.(int)
				if ok {
					floatData = append(floatData, float64(num))
					continue
				}
				numFloat64, ok := v.(float64)
				if ok {
					floatData = append(floatData, numFloat64)
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
				num, ok := v.(int)
				if ok {
					floatData = append(floatData, float64(num))
					continue
				}
				numFloat64, ok := v.(float64)
				if ok {
					floatData = append(floatData, numFloat64)
					continue
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

		fmt.Println(binaryRatio, ratioAmount)
		var treeNodes []*TreeStructure
		root := &TreeStructure{UserID: 1, Levels: 0, Cycle: 1}
		result := buildTree(root, int(numMembers), 2, floatData)
		for _, tree := range result {
			treeNodes = append(treeNodes, tree)
		}
		stored_id, cycleCount, all_data := AllocateMembers(int(numMembers), intData, 2, result)
		AssignJoiningFee(treeNodes, stored_id, floatData)

		totalSponsorBonus := CalculateSponsorBonus(treeNodes, sponsorPercentage, floatData, cappingAmount, cappingScope, int(numMembers))
		// totalBinaryBonus := BinaryWithRatio(all_data, floatData, binaryRatio, int(ratioAmount), cappingScope, cappingAmount)
		totalBinaryBonus := CalculateBinaryBonus(treeNodes, binaryPercentage, cappingAmount, cappingScope)
		var temp []TreeStructureJSON
		var all_data_json [][]TreeStructureJSON
		for _, list := range all_data {
			temp = convertToJSONStructure(list)
			all_data_json = append(all_data_json, temp)
		}
		totalMatchingBonus := CalculateMatchingBonus(treeNodes, floatData, cappingAmount, cappingScope)
		results := map[string]interface{}{
			"tree_structure": convertToJSONStructure(treeNodes),
			"stored_id":      stored_id,
			"cycleCount":     cycleCount,
			"all_data":       all_data_json,
			"total_sponsor_bonus": totalSponsorBonus,
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
*/
