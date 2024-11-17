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
}
type TreeStructureJSON struct {
	UserID              int     `json:"user_id"`
	Position            string  `json:"position"`
	Joining_package_fee float64 `json:"joining_package_fee"`
	Levels              int     `json:"levels"`
	Lft                 int     `json:"lft"`
	Rgt                 int     `json:"rgt"`
	ParentID            int     `json:"parentid"`
	SponsorBonus        float64 `json:"sponsor_bonus"`
	BinaryBonus         float64 `json:"binary_bonus"`
	MatchingBonus       float64 `json:"matching_bonus"`
	CarryForward        float64 `json:"carry_forward"`
	Cycle               int     `json:"cycle"`
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
	if numMembers <= 0{
		fmt.Println("Invalid input: numMembers or cycles is non-positive")
		return nil
	}

	treeNodes := []*TreeStructure{root}
	queue := []*TreeStructure{root}
	currentID := startID
	// currentCycle := 1
	nodesToAdd := numMembers
	for len(queue) > 0 && nodesToAdd > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.LeftMember == nil && nodesToAdd > 0 {
			leftChild := &TreeStructure{
				UserID:              currentID,
				ParentID:            node,
				Levels:              node.Levels + 1,
				Position:            "Left",
				// Cycle:               currentCycle,
			}
			node.LeftMember = leftChild
			queue = append(queue, leftChild)
			treeNodes = append(treeNodes, leftChild)
			currentID++
			nodesToAdd--

			if nodesToAdd > 0 {
				rightChild := &TreeStructure{
					UserID:              currentID,
					ParentID:            node,
					Levels:              node.Levels + 1,
					Position:            "Right",
					// Cycle:               currentCycle,
				}
				node.RightMember = rightChild
				queue = append(queue, rightChild)
				treeNodes = append(treeNodes, rightChild)
				currentID++
				nodesToAdd--
			}
		} else {
			rightChild := &TreeStructure{
				UserID:              currentID,
				ParentID:            node,
				Levels:              node.Levels + 1,
				Position:            "Right",
				// Cycle:               currentCycle,
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

func AllocateMembers(numMembers int, product_quantity []int, startID int) [][]int {
	currentID := startID
	remaing := numMembers + 2
	var allCycles[][]int
	for remaing > 2{
		var cycle[][]int
		for _, qty := range product_quantity {
			cycleID := []int{}
			for j := 0; j < qty && remaing > 2; j++ {
				cycleID = append(cycleID, currentID)
				currentID++ 
				remaing--
			}
			cycle = append(cycle, cycleID)
			if remaing <= 0{
				break
			}
		}
		for _, lst := range cycle {
			if len(lst) > 0 {
				allCycles = append(allCycles, lst)
			}
		}
	}		
	fmt.Println(allCycles)
	return allCycles
}
	


// func AssignJoiningFee(nodes []*TreeStructure,joining_package_fee []float64, numMembers int) float64{
// 	allExceptFirst := nodes[1:]
// 	var totalJoiningFee float64
//     currentCycle := 1
//     var cycleMemberCount int
//     for _, member := range allExceptFirst {      
//         if currentCycle <= len(joining_package_fee) {
//             member.Joining_package_fee = joining_package_fee[currentCycle-1]
// 			totalJoiningFee += member.Joining_package_fee
//         }      
//         cycleMemberCount++
//         if cycleMemberCount == numMembers {
//             cycleMemberCount = 0
//             currentCycle++
//         }
//     }
// 	return totalJoiningFee
// }

func AssignJoiningFee(nodes []*TreeStructure, ids [][]int, joining_package_fee []float64){
	allExceptFirst := nodes[1:]
	var jpf float64
	for i :=range ids{
		if i<len(joining_package_fee){
			jpf = joining_package_fee[i]
		}else{
			jpf = joining_package_fee[i - len(joining_package_fee)]
		}
		for _, id := range ids[i]{
			for _, node := range allExceptFirst{
				if id == node.UserID{
					node.Joining_package_fee = jpf
					break
				}
			}
		}
	}
}

func CalculateSponsorBonus(nodes []*TreeStructure,sponsorBonusPercent float64,joining_package_fee []float64, cappingAmount float64,cappingScope string, numMembers int,) float64 {
    var totalBonus float64
    currentCycle := 1
    var cycleMemberCount int

    for _, member := range nodes {
        var rightBonus, leftBonus float64

        
        if currentCycle <= len(joining_package_fee) {
            member.Joining_package_fee = joining_package_fee[currentCycle-1]
        }

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

func Traverse(node *TreeStructure, joining_package_fee float64) float64 {
	if node == nil {
		return 0
	}
	currentSales := joining_package_fee

	leftSales := Traverse(node.LeftMember, joining_package_fee)
	rightSales := Traverse(node.RightMember, joining_package_fee)

	return currentSales + leftSales + rightSales
}
func CalculateBinaryBonus(nodes []*TreeStructure, binaryBonusPercentage, cappingAmount float64, cappingScope string, joining_package_fee float64) float64 {
	var totalBonus float64

	for _, member := range nodes {
		leftSales := 0.0
		rightSales := 0.0

		if member.LeftMember != nil {
			leftSales = Traverse(member.LeftMember, joining_package_fee)
			member.LeftSales = leftSales
		}
		if member.RightMember != nil {
			rightSales = Traverse(member.RightMember, joining_package_fee)
			member.RightSales = rightSales
		}

		binaryBonus := math.Min(leftSales, rightSales) * binaryBonusPercentage / 100

		if strings.Contains(cappingScope, "binary") && binaryBonus > cappingAmount {
			member.BinaryBonus = cappingAmount
		} else {
			member.BinaryBonus = binaryBonus
		}

		carryForward := leftSales - rightSales
		if member.LeftMember != nil && carryForward > 0 {
			member.LeftMember.CarryForward = carryForward
		} else if member.RightMember != nil && carryForward < 0 {
			member.RightMember.CarryForward = -carryForward
		}
		totalBonus += member.BinaryBonus
	}

	return totalBonus
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
		// cycles, ok := data["cycle"].(float64)
		// if !ok {
		// 	http.Error(w, "Invalid or missing 'cycle' field", http.StatusBadRequest)
		// 	return
		// }
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
		// var joiningPackageFeeData []float64
		// for _, v := range joining_package_fee {
		// 	if num, ok := v.(int); ok {
		// 		joiningPackageFeeData = append(joiningPackageFeeData, float64(num))
		// 	}
		// }
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
		// var productQuantityData []float64
		// for _, v := range product_quantity {
		// 	if num, ok := v.(int); ok {
		// 		productQuantityData = append(productQuantityData, float64(num))
		// 	}
		// }
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

		var treeNodes []*TreeStructure
		root := &TreeStructure{UserID: 1, Levels: 0, Cycle: 1}
		result := buildTree(root, int(numMembers), 2, floatData)
		for _, tree := range result {
			treeNodes = append(treeNodes, tree)
		}
		stored_id := AllocateMembers(int(numMembers), intData, 2)
		// totalJoiningBonus := AssignJoiningFee(treeNodes,floatData, int(numMembers)) 
		AssignJoiningFee(treeNodes, stored_id, floatData)
		for _,node := range(treeNodes){
			fmt.Println(node.UserID, node.Joining_package_fee)
		}
		totalSponsorBonus := CalculateSponsorBonus(treeNodes, sponsorPercentage, floatData, cappingAmount, cappingScope,int(numMembers))
		// totalSponsorBonus = CalculateSponsorBonus(treeNodes, sponsorPercentage, floatData, cappingAmount, cappingScope,int(numMembers))
		totalMatchingBonus := CalculateMatchingBonus(treeNodes, floatData, cappingAmount, cappingScope)
		results := map[string]interface{}{
			"tree_structure":      convertToJSONStructure(treeNodes),
			"stored_id": stored_id,
			//"total_joining_bonus": totalJoiningBonus,
			"total_sponsor_bonus": totalSponsorBonus,
			// "total_binary_bonus":   totalBinaryBonus,
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
