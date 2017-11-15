// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario
 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
*/

package main

/* Imports
* 4 utility libraries for handling bytes, reading and writing JSON,
formatting, and string manipulation
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts
*/
import (
    "bytes"
    "encoding/json"
    "fmt"
    "strconv"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Cow and Pork structure, with properties.
Structure tags are used by encoding/json library
*/
type Register_info struct { // *소출생 등 신고정보[]
    Owner         string `json:"owner"`         // 소유주
    Category      string `json:"category"`      // 신고구분
    Date          string `json:"date"`          // 년월일
    Owner_address string `json:"owner_address"` // 사육지(주소)
}

type Slaughter_info struct { // *도축정보
    Slaughter_house string `json:"slaughter_house"` // 도축장
    Slaughter_date  string `json:"slaughter_date"`  // 도축일자
    Cow_result      string `json:"cow_result"`      // 도축검사결과
    Cow_weight      int    `json:"cow_weight"`      // 도체중
    Cow_grade       string `json:"cow_grade"`       // 육질등급
    Slaughter_company string `json:"slaughter_company"` // 도축처리업소
}

type Foot_and_mouth struct { // *구제역
    Vaccine_date   string `json:"vaccine_date"`   // 구제역 예방접종(최>종)일자
    Vaccine_result string `json:"vaccine_result"` // 구제역 검사결과
}

type Brucelliasis struct { // *브루셀라
    Vaccine_date   string `json:"vaccine_date"`   // 브루셀라 검사(최종)일자
    Vaccine_result string `json:"vaccine_result"` // 브루셀라 검사결과
}

type Tuberculosis struct { // *결핵
    Vaccine_date   string `json:"vaccine_date"`   // 결핵 검사(최종)일자
    Vaccine_result string `json:"vaccine_result"` // 결핵 검사결과
}

type Package_info struct { // *포장정보
    Company        string `json:"company"`          // 포장회사
    Company_address string `json:"company_address"` // 포장회사주소
    Cow_part       string `json:"cow_part"`         // 포장부위
    Package_amount string `json:"package_amoount"`  // 포장단위(g)
    Package_date   string `json:"package_date"`     // 포장일자
}

type Cow struct { // *소개체 정보
    Trace_id       string           `json:"trace_id"`       // 이력번호
    Farm_id        string           `json:"farm_id"`        // 농장식별번호
    Cow_id         string           `json:"cow_id"`         // 개체식별번호
    Cow_birthday   string           `json:"cow_birthday"`   // 출생년월일
    Cow_category   string           `json:"cow_category"`   // 소의 종류
    Cow_sex        string           `json:"cow_sex"`        // 성별
    Register_info  []Register_info  `json:"register_info"`  // 소출생등신고정보[]
    Slaughter_info Slaughter_info   `json:"slaughter_info"` // 도축정보
    Foot_and_mouth []Foot_and_mouth `json:"foot_and_mouth"` // 구제역[]
    Brucelliasis   []Brucelliasis   `json:"brucelliasis"`   // 브루셀[]
    Tuberculosis   []Tuberculosis   `json:"tuberculosis"`   // 결핵[]
}

/*
Pork =                돼지개체 정보 = {
id                    이력번호
farm_id               농장식별번호 (??소유주가 바뀌면 농장식별번호는 바뀌나??)
farm_owner            농장경영자
farm_address          농장소재지
slaughterhouse        도축장
slaughter_address     소재지
slaughter_date        도축일자
inspection_result     도축검사결과
package_info = [
  {
     company_name     업소명
     company_address  소재지
     package_date     포장일자
     package_grade    등급
  }
]
}
*/

/*
 * The Init method *
 called when the Smart Contract "traceability-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function
 -- see initLedger()
*/
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
    return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "traceability-chaincode"
 The app also specifies the specific smart contract function to call with args
*/
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

    // Retrieve the requested Smart Contract function and arguments
    function, args := APIstub.GetFunctionAndParameters()
    // Route to the appropriate handler function to interact with the ledger

    if function == "queryCow" {
        return s.queryCow(APIstub, args)
    } else if function == "initLedger" {
        return s.initLedger(APIstub)
    } else if function == "recordCow" {
        return s.recordCow(APIstub, args)
    } else if function == "queryAllCow" {
        return s.queryAllCow(APIstub)
    } else if function == "changeCowHolder" {
        return s.changeCowHolder(APIstub, args)
    } else if function == "registerCow" {
        return s.registerCow(APIstub, args)
    } else if function == "updateSlaughterInfoCow" {
        return s.updateSlaughterInfoCow(APIstub, args)
    } else if function == "updatePackageInfoCow" {
        return s.updatePackageInfoCow(APIstub, args)
    }

    return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryCow method *
Used to view the records of one particular cow
It takes one argument -- the key for the cow in question
*/
func (s *SmartContract) queryCow(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    cowAsBytes, _ := APIstub.GetState(args[0])
    if cowAsBytes == nil {
        return shim.Error("Could not locate cow")
    }
    return shim.Success(cowAsBytes)
}

/*
 * The initLedger method *
Will add test data (10 cow catches)to our network
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
    cows := []Cow{
        Cow{
            Trace_id:     "002021864077",
            Farm_id:      "020218",
            Cow_id:       "64077",
            Cow_birthday: "20070115",
            Cow_category: "한우",
            Cow_sex:      "거세",
            Register_info: []Register_info{
                Register_info{
                    Owner:         "김영한",
                    Category:      "전산등록",
                    Date:          "20080111",
                    Owner_address: "경기도 양주시 만송통",
                },
                Register_info{
                    Owner:         "양영귀",
                    Category:      "양수",
                    Date:          "20090611",
                    Owner_address: "경기도 양주시 만송통",
                },
                Register_info{
                    Owner:         "양영귀",
                    Category:      "도축출하",
                    Date:          "20090611",
                    Owner_address: "경기도 양주시 만송통",
                },
            },
            Slaughter_info: Slaughter_info{
                Slaughter_house: "우림축산 (경기도 동두천시 동두천동)",
                Slaughter_date:  "20090611",
                Cow_result:      "합격",
                Cow_weight:      409,
                Cow_grade:       "1+",
                Slaughter_company: "양주축협가공공장 (경기도 양주시 고읍동)",
            },
            Foot_and_mouth: []Foot_and_mouth{
                Foot_and_mouth{
                    Vaccine_date:   "20080211",
                    Vaccine_result: "이상무",
                },
            },
            Brucelliasis: []Brucelliasis{
                Brucelliasis{
                    Vaccine_date:   "20080420",
                    Vaccine_result: "합격",
                },
            },
            Tuberculosis: []Tuberculosis{
                Tuberculosis{
                    Vaccine_date:   "20080631",
                    Vaccine_result: "합격",
                },
            },
        },
        Cow{
            Trace_id:           "002021864078",
            Farm_id:      "020218",
            Cow_id:       "64077",
            Cow_birthday: "20070115",
            Cow_category: "한우",
            Cow_sex:      "거세",
            Register_info: []Register_info{
                Register_info{
                    Owner:         "김영한",
                    Category:      "전산등록",
                    Date:          "20080111",
                    Owner_address: "경기도 양주시 만송통",
                },
                Register_info{
                    Owner:         "양영귀",
                    Category:      "양수",
                    Date:          "20090611",
                    Owner_address: "경기도 양주시 만송통",
                },
                Register_info{
                    Owner:         "양영귀",
                    Category:      "도축출하",
                    Date:          "20090611",
                    Owner_address: "경기도 양주시 만송통",
                },
            },
            Slaughter_info: Slaughter_info{
                Slaughter_house: "우림축산 (경기도 동두천시 동두천동)",
                Slaughter_date:  "20090611",
                Cow_result:      "합격",
                Cow_weight:      409,
                Cow_grade:       "1+",
                Slaughter_company: "양주축협가공공장 (경기도 양주시 고읍동)",
            },
            Foot_and_mouth: []Foot_and_mouth{
                Foot_and_mouth{
                    Vaccine_date:   "20080211",
                    Vaccine_result: "이상무",
                },
            },
            Brucelliasis: []Brucelliasis{
                Brucelliasis{
                    Vaccine_date:   "20080420",
                    Vaccine_result: "합격",
                },
            },
            Tuberculosis: []Tuberculosis{
                Tuberculosis{
                    Vaccine_date:   "20080631",
                    Vaccine_result: "합격",
                },
            },
        },
    }

    i := 0
    for i < len(cows) {
        fmt.Println("i is ", i)
        cowAsBytes, _ := json.Marshal(cows[i])
        //APIstub.PutState(strconv.Itoa(i+1), cowAsBytes)
        APIstub.PutState(cows[i].Trace_id, cowAsBytes)
        fmt.Println("Added", cows[i])
        i = i + 1
    }

    return shim.Success(nil)
}

/*
 * The recordCow method *
Fisherman like Sarah would use to record each of her cow catches.
This method takes in five arguments (attributes to be saved in the ledger).
*/
func (s *SmartContract) recordCow(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 5 {
        return shim.Error("Incorrect number of arguments. Expecting 5")
    }

    var cow = Cow{Trace_id: args[1], Farm_id: args[2], Cow_id: args[3], Cow_birthday: args[4],
                  Cow_category: args[5], Cow_sex: args[6]}

    cowAsBytes, _ := json.Marshal(cow)
    err := APIstub.PutState(args[0], cowAsBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("Failed to record cow catch: %s", args[0]))
    }

    return shim.Success(nil)
}

/*
 * The queryAllCow method *
allows for assessing all the records added to the ledger(all cow catches)
This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SmartContract) queryAllCow(APIstub shim.ChaincodeStubInterface) sc.Response {

    startKey := "000000000000"
    endKey := "999999999999"

    resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    // buffer is a JSON array containing QueryResults
    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        // Add comma before array members,suppress it for the first array member
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"Key\":")
        buffer.WriteString("\"")
        buffer.WriteString(queryResponse.Key)
        buffer.WriteString("\"")

        buffer.WriteString(", \"Record\":")
        // Record is a JSON object, so we write as-is
        buffer.WriteString(string(queryResponse.Value))
        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")

    fmt.Printf("- queryAllCow:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}

/*
 * The changeCowHolder method *
The data in the world state can be updated with who has possession.
This function takes in 2 arguments, cow id and new holder name.
*/
func (s *SmartContract) changeCowHolder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }

    cowAsBytes, _ := APIstub.GetState(args[0])
    if cowAsBytes == nil {
        return shim.Error("Could not locate cow")
    }
    cow := Cow{}

    json.Unmarshal(cowAsBytes, &cow)
    // Normally check that the specified argument is a valid holder of cow
    // we are skipping this check for this example
    cow.Cow_birthday = args[1]

    cowAsBytes, _ = json.Marshal(cow)
    err := APIstub.PutState(args[0], cowAsBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("Failed to change cow holder: %s", args[0]))
    }

    return shim.Success(nil)
}

/*
 * The registerCow method *
*/
func (s *SmartContract) registerCow(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 8 {
        return shim.Error("Incorrect number of arguments. Expecting 5")
    }
    
    if len(args[0]) != 12 {
        return shim.Error("Inccorrect trace id. Expecting 12 digits")
    }

    var register_info []Register_info
    var farm_id string = args[0][1:7] 
    var cow_id string = args[0][7:]

    register_info = append(register_info, Register_info{ Owner: args[4], Category: args[5], Date: args[6], Owner_address: args[7] })

    var cow = Cow{ Trace_id: args[0], Farm_id: farm_id, Cow_id: cow_id, Cow_birthday: args[1],
        Cow_category: args[2], Cow_sex: args[3],
        //Register_info: Register_info[0]{ Owner: args[6], Category: args[7], Date: args[8], Owner_address: args[9] }}
        Register_info: register_info}

    cowAsBytes, _ := json.Marshal(cow)
    err := APIstub.PutState(args[0], cowAsBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("Failed to register: %s", args[0]))
    }

    return shim.Success(nil)
}

/*
 * The updateSlaughterInfoCow method *
*/
func (s *SmartContract) updateSlaughterInfoCow(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 7 {
        return shim.Error("[updateSlaughterInfoCow] Incorrect number of arguments. Expecting 7")
    }

    if len(args[0]) != 12 {
        return shim.Error("Inccorrect trace id. Expecting 12 digits")
    }

    weight, err := strconv.Atoi(args[4])
    if err != nil {
        // handle error
        return shim.Error("Incorrect number for cow weight")
    }

    cowAsBytes, _ := APIstub.GetState(args[0])
    if cowAsBytes == nil {
        return shim.Error("Could not locate cow")
    }

    var cow = Cow{}

    json.Unmarshal(cowAsBytes, &cow)

    var slaughter_info Slaughter_info
    slaughter_info = Slaughter_info{ Slaughter_house: args[1], Slaughter_date: args[2], Cow_result: args[3],
                                     Cow_weight: weight, Cow_grade: args[5], Slaughter_company: args[6] }

    cow.Slaughter_info = slaughter_info

    cowAsBytes, _ = json.Marshal(cow)
    err = APIstub.PutState(args[0], cowAsBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("Failed to register slaughter information: %s", args[0]))
    }

    return shim.Success(nil)
}

/*
 * The updatePackageInfoCow method *
*/
func (s *SmartContract) updatePackageInfoCow(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 6 {
        return shim.Error("[updatePackageInfoCow] Incorrect number of arguments. Expecting 6")
    }

    if len(args[0]) != 12 {
        return shim.Error("Inccorrect trace id. Expecting 12 digits")
    }

    package_amount, err := strconv.Atoi(args[4])
    if err != nil {
        // handle error
        return shim.Error("Incorrect number for package amount")
    }

    cowAsBytes, _ := APIstub.GetState(args[0])
    if cowAsBytes == nil {
        return shim.Error("Could not locate cow")
    }

    var cow = Cow{}

    json.Unmarshal(cowAsBytes, &cow)

    var package_info Package_info
    package_info = Package_info{ Company: args[1], Company_address: args[2], Cow_part: args[3],
                                 Package_amount: package_amount, Package_date: args[5] }

    fmt.Print("PackageInfo")
    cow.Package_info = append(cow.Package_info, package_info)

    cowAsBytes, _ = json.Marshal(cow)
    err = APIstub.PutState(args[0], cowAsBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("Failed to register slaughter information: %s", args[0]))
    }

    return shim.Success(nil)
}

/*
 * main function *
calls the Start functon
The main function starts the chaincode in the container during instantiation.
*/
func main() {

    // Create a new Smart Contract
    err := shim.Start(new(SmartContract))
    if err != nil {
        fmt.Printf("Error creating new Smart Contract: %s", err)
    }
}
