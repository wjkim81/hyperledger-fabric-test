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

/* Define Cattle and Pig structure, with properties.
Structure tags are used by encoding/json library
*/
/*
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
    Package_amount int    `json:"package_amount"`  // 포장단위(g)
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
    Package_info   []Package_info   `json:"package_info"`   // 포장정보[]
    Foot_and_mouth []Foot_and_mouth `json:"foot_and_mouth"` // 구제역[]
    Brucelliasis   []Brucelliasis   `json:"brucelliasis"`   // 브루셀[]
    Tuberculosis   []Tuberculosis   `json:"tuberculosis"`   // 결핵[]
}
*/
/*
data.go.kr에서 제공하는 open-api 데이터 조회를 기준으로 작성
- 입력값
TraceNo: 개체/이력/묶음번호
OptionNo: 조회하려는 정보 옵션값
    1~7: 개체/이력, 8~9 : 묶음
    1: 개체(소), 사육(돼지)
    2: 출생 등 신고(소)
    3: 도축(소/돼지)
    4: 포장(소/돼지)
    5: 구제역백신(소)
    6: 질병정보(소)
    7: 브루셀라(소)
    8: 묶음 기본정보(묶음)
    9: 묶음구성내역(묶음)
CorpNo:  묶음을 구성한 업소의 사업자 번호

- 출력값
기본 정보
traceNoType: 소/돼지/묶음 구분
    CATTLE/CATTLE_NO: 소 개체
    CATTLE/LOT_NO: 소 묶음
    PIG/PIG_NO: 돼지 이력
    PIG/LOT_NO: 돼지 묶음
infoType: 조회 정보 옵션값
    1~7: 개체/이력, 8~9 : 묶음
    1: 개체(소), 사육(돼지)
    2: 출생 등 신고(소)
    3: 도축(소/돼지)
    4: 포장(소/돼지)
    5: 구제역백신(소)
    6: 질병정보(소)
    7: 브루셀라(소)
    8: 묶음 기본정보(묶음)
    9: 묶음구성내역(묶음)

    항목명(영문)     항목명(국문)   항목크기    항목구분    샘플데이터         항목설명
traceNoType: CATTLE/CATTLE_NO: 소 개체
infoType = 1 o
    BirthYmd      출생일자      10        1         20120317        출생년월일
    CattleNo      소 개체번호    15       1         410002075264204 이력번호
    FlatEartagNo  재귀표번호     10       0          x               재부착번호
    LsTypeNm      소의 종류     20       1          한우             소의 종류
    MonthDiff     수입경과월     10       0          0              수입경과월
    NationNm      수입국가      20       0          x               수입국가
    SexNm         성별         10       1          암              성별
    FarmUniqueNo  농장식별번호   6        0          211058          농장식별번호

infoType = 2 o
    FarmAddr      사육지       200      1          강원도 원주시 호저면  사육지
    FarmerNm      소유주       100      1          고정남             소유주
    RegType       신고구분     20        1         전산등록            신고구분
    RegYmd        년월일       8        1          20120319         년월일

infoType = 3
    ButcheryPlaceAddr  도축장 주소   200 1          경기도 안성시 일죽면 일생로  도축장 주소
    ButcheryPlaceNm    도축장명     100 1           도드람LPC              도축장명
    ButcheryYmd        도축일자     8   1           20170920             도축일자
    GradeNm            등급명      20  0           1                    육질등급
    InspectPassYn      위생검사 결과 1   0           합격                  위생검사결과 합격/불합격/보류

infoType = 4
    ProcessPlaceAddr   포장처리업소 주소   200    1   경기도 광주시 곤지암읍 경충대로   포장처리업소 주소
    ProcessPlaceNm     포장처리업소명     50     1   우성육가공                    포장처리업소명

infoType = 5
    InjectionDayCnt 구제역 백신접종경과일   200    0       접종 후 2일 경과       구제역 백신접종경과일
    InjectionYmd    구제역 백신접종일     8       0       20170918           구제역 예방접종최종일자
    Vaccineorder    구제역 백신접종 차수   10      0       13차 구제역          백신접종 차수

infoType = 6
    InspectDesc      질병유무          100     0        x                   질병유무

infoType = 7
    InspectDt       브루셀라 접종일      8       0       20170911        브루셀라 검사최종일자
    InspectYn       브루셀러 접종 유무    10      0       음성             브루셀라 검사결과

traceNoType: CATTLE/LOT_NO: 소묶음
infoType = 8
    CorpNo           사업자번호        10       1       1178522046                사업자번호
    LotNo            묶음번호         30       1        L01709271277007          묶음번호
    ProcessPlaceAddr 묶음구성업소 주소   200     1       서울특별시 양천구 목동동로 257   묶음구성업소 주소
    ProcessPlaceNm   묶음구성업소명     50       1       ㈜현대그린푸드목동점           묶음구성업소명

infoType = 9
    CattleNo            소 개체번호  15  1   410002075264204                 소 개체번호
    CorpNo              사업자번호   10  1   1178522046                      사업자번호
    FarmAddr            농장주소    200 1   강원도 원주시 호저면 매호리            농장주소
    GradeNm             등급명     20  1   1                                등급명
    LotNo               묶음번호    30  1   L01709271277007                  묶음번호
    LsTypeNm            소의 종류   20  1   한우  소의 종류
    ProcessPlaceAddr    묶음구성업소 주소   200 1   서울특별시 양천구 목동 916번지  포장처리업소 주소
    ProcessPlaceNm      묶음구성업소명 50  1   ㈜현대그린푸드목동점               묶음구성업소명
    ButcheryPlaceAddr   도축장 주소  200 1   경기도 안성시 일죽면 금산리 598번지   도축장 주소
    ButcheryPlaceNm     도축장명    100 1   도드람LPC                       도축장명
    ButcheryYmd         도축일자    8   1   20170920                       도축일자



traceNoType: PIG/PIG_NO: 돼지 개체
infoType = 1
    FarmAddr             농장소재지   200 1   경상북도 포항시 북구 청하면 비학로     농장소재지
    FarmerNm             농장경영자   100 1   강충열                          농장경영자
    FarmUniqueNo         농장식별번호  6   1   700030                        농장식별번호
    PigNo                이력번호    20  1   170003000058                   이력번호

infoType = 2
    ButcheryPlaceAddr    도축장 주소  200 1   경상북도 영천시 동강포길 29-5  도축장 주소
    ButcheryPlaceNm      도축장명    100 1   (주)삼세   도축장명
    ButcheryYmd          도축일자    8   1   20150730    도축일자
    InspectPassYn        위생검사 결과     0   Y   합격 : Y, 불합격 : N, 보류 : H
    PigNo                돼지 이력번호 20  1   170003000058    돼지 이력번호

infoType = 3
    GradeNm              등급명      20  0    -                             등급명
    PigNo                돼지 이력번호 20  1   170003000058                   돼지 이력번호
    ProcessPlaceAddr     포장처리업소 주소  200 1   경상남도 양산시 상북면 ,지하1층   포장처리업소 주소
    ButcheryPlaceNm      도축장명     100 1      ㈜산청푸드                    도축장명
    ProcessYmd           포장처리일자   8   1   20150731                     포장처리일자

infoType = 8
    CorpNo               사업자번호   10  1   6218176868                    사업자번호
    GradeNm              등급명      20  0                                 등급명
    LotNo                묶음번호    30  1    L11507319120001               묶음번호
    ProcessPlaceAddr     포장처리업소 소재지  200 1   경상남도 양산시 상북면 충렬로 574, 지하1층  포장처리업소 소재지
    ProcessPlaceNm       포장처리업소명 50  1   ㈜산청푸드                      포장처리업소명
    ProcessYmd           포장처리일자  8   1   20150731                      포장처리일자

infoType = 9
    CorpNo               사업자번호   10  1   1992484871                    사업자번호
    FarmAddr             농장 소재지  200 1   경상북도 포항시 북구 청하면 비학로 농장 소재지
    FarmerNm             사육자명    100 1   강충열                           사육자명
    PigNo                이력(묶음)번호    20  1   170003000058              이력(묶음)번호
    ButcheryPlaceAddr    도축장 주소  200 1   경상북도 영천시 도남동 695번지      도축장 주소
    ButcheryPlaceNm      도축장명    100 1   (주)삼세                       도축장명
    ButcheryYmd          도축일자    8   1   20150730                     도축일자
*/

type FarmInfo struct { // 농장 정보 및 등록정보
	FarmNo   string `json:"farmNo"`   // 농장식별번호
	FarmNm   string `json:"farmNm"`   // 농장명
	FarmAddr string `json:"farmAddr"` // 농장주소
	FarmerNm string `json:"farmerNm"` // 소유주
	RegType  string `json:"regType"`  // 신고구분
	RegYmd   string `json:"regYmd"`   // 년월일
}

type ButcheryInfo struct {
	ButcheryPlaceAddr string `json:"butcheryPlaceAddr"` // 도축장 주소
	ButcheryPlaceNm   string `json:"butcheryPlaceNm"`   // 도축장명
	ButcheryYmd       string `json:"butcheryYmd"`       // 도축일자
	GradeNm           string `json:"gradeNm"`           // 등급명 (육질등급)
	InspectPassYn     string `json:"inspectPassYn"`     // 위생검사 결과 (합격/불합격/보류)
	ButcheryWeight    int    `json:"butcheryWeight"`    // 도체중
	AbattCode         string `json:"AbattCode"`         // 도축장코드
	ProcessPlaceNm    string `json:"processPlaceNm"`    // 가공장/도축처리업소
}

type FootAndMouth struct { // *구제역
	InjectionYmd    string `json:"injectionYmd"`    // 구제역 예방접종일자
	InjectionResult string `json:"injectionResult"` // 구제역 검사결과
}

type Brucelliasis struct { // *브루셀라
	InjectionYmd    string `json:"injectionYmd"`    // 브루셀라 검사(최종)일자
	InjectionResult string `json:"injectionResult"` // 브루셀라 검사결과
}

type Tuberculosis struct { // *결핵
	InjectionYmd    string `json:"injectionYmd"`    // 결핵 검사(최종)일자
	InjectionResult string `json:"injectionResult"` // 결핵 검사결과
}

type ProcessInfo struct {
	CorpNo           string `json:"corpNo"`           // 사업자번호
	LotNo            string `json:"lotNo"`            // 묶음번호
	ProcessPlaceNm   string `json:"processPlaceNm"`   // 포장처리업소명
	ProcessPlaceAddr string `json:"processPlaceAddr"` // 포장처리업소 주소
	ProcessYmd       string `json:"processYmd"`       // 포장일자
	ProcessWeight    int    `json:"processWeight"`    // 포장단위(g)
	ProcessPart      string `json:"processPart"`      // 포장 부위
}

type Cattle struct {
	TraceId      string `json:"traceId"`      // 이력번호
	CattleNo     string `json:"cattleNo"`     // 개체식별번호 = 이력번호
	BirthYmd     string `json:"birthYmd"`     // 출생년월일
	FlatEartagNo string `json:"flatEartagNo"` // 재부착번호
	LsTypeCd     string `json:"lsTypeCd"`     // 소의종류코드
	LsTypeNm     string `json:"lsTypeNm"`     // 소의종류
	MonthDiff    int    `json:"monthDiff"`    // 수입경과월
	nationNm     string `json:"nationNm"`     // 수입국가
	SexCd        string `json:"sexCd"`        // 성별코드
	SexNm        string `json:"sexNm"`        // 성별

	FarmInfo     []FarmInfo   `json:"farmInfo"`     // 농장 및 등록정보[]
	ButcheryInfo ButcheryInfo `json:"butcheryInfo"` // 도축정보

	FootAndMouth []FootAndMouth `json:"footAndMouth"` // 구제역[]
	Brucelliasis []Brucelliasis `json:"brucelliasis"` // 브루셀라[]
	Tuberculosis []Tuberculosis `json:"tuberculosis"` // 결핵[]

	ProcessInfo []ProcessInfo `json:"processInfo"` // 포장정보[]
}

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

	if function == "queryCattle" {
		return s.queryCattle(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAllCattle" {
		return s.queryAllCattle(APIstub)
	} else if function == "registerCattle" {
		return s.registerCattle(APIstub, args)
	} else if function == "updateButcheryInfo" {
		return s.updateButcheryInfo(APIstub, args)
	} else if function == "updateProcessInfo" {
		return s.updateProcessInfo(APIstub, args)
	} else if function == "insertObjects" {
		return s.insertObjects(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryCattle method *
Used to view the records of one particular cattle
It takes one argument -- the key for the cattle in question
*/
func (s *SmartContract) queryCattle(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	cattleAsBytes, _ := APIstub.GetState(args[0])
	if cattleAsBytes == nil {
		return shim.Error("Could not locate cattle")
	}
	return shim.Success(cattleAsBytes)
}

/*
 * The initLedger method *
Will add test data to our network
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	cattles := []Cattle{
		Cattle{
			TraceId:  "002021864077",
			CattleNo: "002021864077",
			BirthYmd: "20070115",
			//FlatEartagNo:    "",
			//LsTypeCd:        "",
			LsTypeNm: "한우",
			//SexCd:           "",
			SexNm: "거세",

			FarmInfo: []FarmInfo{
				FarmInfo{
					FarmNo:   "020218",
					FarmNm:   "",
					FarmerNm: "김영한",
					FarmAddr: "경기도 양주시 만송통",
					RegType:  "전산등록",
					RegYmd:   "20080111",
				},
				FarmInfo{
					FarmNo:   "020219",
					FarmNm:   "",
					FarmerNm: "양영귀",
					FarmAddr: "경기도 양주시 만송통",
					RegType:  "양수",
					RegYmd:   "20090611",
				},
				FarmInfo{
					FarmNo:   "020219",
					FarmNm:   "",
					FarmerNm: "양영귀",
					FarmAddr: "경기도 양주시 만송통",
					RegType:  "도축출하",
					RegYmd:   "20090611",
				},
			},

			ButcheryInfo: ButcheryInfo{
				ButcheryPlaceAddr: "경기도 동두천시 동두천동",
				ButcheryPlaceNm:   "우림축산",
				ButcheryYmd:       "20090611",
				GradeNm:           "1+",
				InspectPassYn:     "합격",
				ButcheryWeight:    409,
				//AbattCode:    "",
				ProcessPlaceNm: "양주축협가공공장 (경기도 양주시 고읍동)",
			},
			FootAndMouth: []FootAndMouth{
				FootAndMouth{
					InjectionYmd:    "20080211",
					InjectionResult: "이상무",
				},
			},
			Brucelliasis: []Brucelliasis{
				Brucelliasis{
					InjectionYmd:    "20080420",
					InjectionResult: "합격",
				},
			},
			Tuberculosis: []Tuberculosis{
				Tuberculosis{
					InjectionYmd:    "20080631",
					InjectionResult: "합격",
				},
			},
			ProcessInfo: []ProcessInfo{
				ProcessInfo{
					CorpNo:           "1178522046",
					LotNo:            "L01709271277007",
					ProcessPlaceAddr: "경기도 광주시 곤지암읍 경충대로",
					ProcessPlaceNm:   "우성 육가공",
					ProcessYmd:       "20190304",
					ProcessWeight:    1000,
					ProcessPart:      "안심",
				},
				ProcessInfo{
					CorpNo:           "1178522046",
					LotNo:            "L01709271277007",
					ProcessPlaceAddr: "경기도 광주시 곤지암읍 경충대로",
					ProcessPlaceNm:   "우성 육가공",
					ProcessYmd:       "20190304",
					ProcessWeight:    2000,
					ProcessPart:      "등심",
				},
			},
		},
		/*
		   Cow{
		       TraceId:           "002021864078",
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
		*/
	}

	i := 0
	for i < len(cattles) {
		fmt.Println("i is ", i)
		cattleAsBytes, _ := json.Marshal(cattles[i])
		//APIstub.PutState(strconv.Itoa(i+1), cowAsBytes)
		APIstub.PutState(cattles[i].TraceId, cattleAsBytes)
		fmt.Println("Added", cattles[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The initLedger method *
Will add test data (10 cattle catches)to our network
*/
func (s *SmartContract) insertObjects(APIstub shim.ChaincodeStubInterface) sc.Response {

	var cattle Cattle
	var traceId int
	var traceIdStr string
	var cattleNoStr string

	cattle = Cattle{
		TraceId:  "002021864077",
		CattleNo: "002021864077",
		BirthYmd: "20070115",
		//FlatEartagNo:    "",
		//LsTypeCd:        "",
		LsTypeNm: "한우",
		//SexCd:           "",
		SexNm: "거세",

		FarmInfo: []FarmInfo{
			FarmInfo{
				FarmNo:   "020218",
				FarmNm:   "",
				FarmerNm: "김영한",
				FarmAddr: "경기도 양주시 만송통",
				RegType:  "전산등록",
				RegYmd:   "20080111",
			},
			FarmInfo{
				FarmNo:   "020219",
				FarmNm:   "",
				FarmerNm: "양영귀",
				FarmAddr: "경기도 양주시 만송통",
				RegType:  "양수",
				RegYmd:   "20090611",
			},
			FarmInfo{
				FarmNo:   "020219",
				FarmNm:   "",
				FarmerNm: "양영귀",
				FarmAddr: "경기도 양주시 만송통",
				RegType:  "도축출하",
				RegYmd:   "20090611",
			},
		},

		ButcheryInfo: ButcheryInfo{
			ButcheryPlaceAddr: "경기도 동두천시 동두천동",
			ButcheryPlaceNm:   "우림축산",
			ButcheryYmd:       "20090611",
			GradeNm:           "1+",
			InspectPassYn:     "합격",
			ButcheryWeight:    409,
			//AbattCode:    "",
			ProcessPlaceNm: "양주축협가공공장 (경기도 양주시 고읍동)",
		},
		FootAndMouth: []FootAndMouth{
			FootAndMouth{
				InjectionYmd:    "20080211",
				InjectionResult: "이상무",
			},
		},
		Brucelliasis: []Brucelliasis{
			Brucelliasis{
				InjectionYmd:    "20080420",
				InjectionResult: "합격",
			},
		},
		Tuberculosis: []Tuberculosis{
			Tuberculosis{
				InjectionYmd:    "20080631",
				InjectionResult: "합격",
			},
		},
		ProcessInfo: []ProcessInfo{
			ProcessInfo{
				CorpNo:           "1178522046",
				LotNo:            "L01709271277007",
				ProcessPlaceAddr: "경기도 광주시 곤지암읍 경충대로",
				ProcessPlaceNm:   "우성 육가공",
				ProcessYmd:       "20190304",
				ProcessWeight:    2000,
				ProcessPart:      "삼겹살",
			},
		},
	}

	traceId = 102021860000
	numInsert := 1000
	for i := 0; i < numInsert; i++ {
		fmt.Println("i is ", i)
		traceId = traceId + i
		traceIdStr = strconv.Itoa(traceId)
		cattleNoStr = strconv.Itoa(100000 + i)

		cattle.TraceId = traceIdStr
		cattle.CattleNo = cattleNoStr
		cattle.LsTypeNm = "돼지 " + strconv.Itoa(i)

		cattleAsBytes, _ := json.Marshal(cattle)
		//APIstub.PutState(strconv.Itoa(i+1), cowAsBytes)
		APIstub.PutState(cattle.TraceId, cattleAsBytes)
		fmt.Println("Added", cattle)
	}

	fmt.Println("Inserting object")

	return shim.Success(nil)
}

/*
 * The queryAllCattle method *
allows for assessing all the records added to the ledger(all cattles)
This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SmartContract) queryAllCattle(APIstub shim.ChaincodeStubInterface) sc.Response {

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

	fmt.Printf("- queryAllCattle:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The registerCattle method *
 */
func (s *SmartContract) registerCattle(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	if len(args[0]) != 12 {
		return shim.Error("Inccorrect trace id. Expecting 12 digits")
	}

	var farmInfo []FarmInfo
	var farmNo string = args[0][1:7]
	var cattleNo string = args[0][7:]

	farmInfo = append(farmInfo, FarmInfo{FarmNo: farmNo, FarmerNm: args[4], RegType: args[5], RegYmd: args[6], FarmAddr: args[7]})

	var cattle = Cattle{TraceId: args[0], CattleNo: cattleNo, BirthYmd: args[1], LsTypeNm: args[2], SexNm: args[3],
		FarmInfo: farmInfo}

	cattleAsBytes, _ := json.Marshal(cattle)
	err := APIstub.PutState(args[0], cattleAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to register: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The updateButcheryInfo method *
 */
func (s *SmartContract) updateButcheryInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 7 {
		return shim.Error("[updateButcheryInfo] Incorrect number of arguments. Expecting 7")
	}

	if len(args[0]) != 12 {
		return shim.Error("Inccorrect trace id. Expecting 12 digits")
	}

	weight, err := strconv.Atoi(args[4])
	if err != nil {
		// handle error
		return shim.Error("Incorrect number for cattle weight")
	}

	cattleAsBytes, _ := APIstub.GetState(args[0])
	if cattleAsBytes == nil {
		return shim.Error("Could not locate cattle")
	}

	var cattle = Cattle{}

	json.Unmarshal(cattleAsBytes, &cattle)

	var butcheryInfo ButcheryInfo
	butcheryInfo = ButcheryInfo{ButcheryPlaceNm: args[1], ButcheryYmd: args[2], InspectPassYn: args[3],
		ButcheryWeight: weight, GradeNm: args[5], ProcessPlaceNm: args[6]}

	cattle.ButcheryInfo = butcheryInfo

	cattleAsBytes, _ = json.Marshal(cattle)
	err = APIstub.PutState(args[0], cattleAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to update butchery information: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The updateProcessInfo method *
 */
func (s *SmartContract) updateProcessInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("[updatePackageInfoCow] Incorrect number of arguments. Expecting 6")
	}

	if len(args[0]) != 12 {
		return shim.Error("Inccorrect trace id. Expecting 12 digits")
	}

	packageAmount, err := strconv.Atoi(args[4])
	if err != nil {
		// handle error
		return shim.Error("Incorrect number for package amount")
	}

	cattleAsBytes, _ := APIstub.GetState(args[0])
	if cattleAsBytes == nil {
		return shim.Error("Could not locate cattle")
	}

	var cattle = Cattle{}

	json.Unmarshal(cattleAsBytes, &cattle)

	var processInfo ProcessInfo
	processInfo = ProcessInfo{ProcessPlaceNm: args[1], ProcessPlaceAddr: args[2], ProcessPart: args[3],
		ProcessWeight: packageAmount, ProcessYmd: args[5]}

	cattle.ProcessInfo = append(cattle.ProcessInfo, processInfo)

	cattleAsBytes, _ = json.Marshal(cattle)
	err = APIstub.PutState(args[0], cattleAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to update process information: %s", args[0]))
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
