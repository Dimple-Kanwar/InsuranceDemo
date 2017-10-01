package main
import ("errors"
        "fmt"
        "strconv"
        "encoding/json"
        "github.com/hyperledger/fabric/core/chaincode/shim")

// Insurance implements a simple chaincode to manage an insurance
type Insurance struct {
}

var policyIndexStr = "_policyIndex" //name for the key/value that will store a list of all known policyNumbers

type Policy struct {
    PolicyNumber string `json:"policyNumber"`
    VehicleNumber string `json:"vehicleNumber"`
    CarUser string `json:"carUser"`
    PremiumNumber string `json:"premiumNumber"`
    CarYear string `json:"carYear"`
    CarModel string `json:"carModel"`
    InsuranceToken string `json:"insuranceToken"`
    Status string `json:"status"`
}


// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func(t * Insurance) Init(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {
    var msg string
    // Get the args from the transaction proposal
    if len(args) != 1 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting ' ' as an argument\", \"code\" : \"503\"}"
        err := stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // Initialize the chaincode
    msg = args[0]
    // Set up any insurance here by calling stub.PutState()
    err := stub.PutState("abc", [] byte(msg))
    if err != nil {
        return nil, err
    }
    // We store the policy Number and the value on the ledger
    var empty[] string
    jsonAsBytes, _:= json.Marshal(empty) //marshal an emtpy array of strings to clear the index
    err = stub.PutState(policyIndexStr, jsonAsBytes)
    if err != nil {
        return nil, err
    }

    fmt.Println("Insurance chaincode is deployed successfully.");
    tosend:= "{ \"message\" : \"Insurance chaincode is deployed successfully.\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", [] byte(tosend))
    if err != nil {
        return nil, err
    }
    return nil, nil
}
// Run - Our entry Dealint for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================

func(t * Insurance) Run(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {
    fmt.Println("run is running " + function)
    return t.Invoke(stub, function, args)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *Insurance) Invoke(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {
    // Extract the function and args from the transaction proposal

    if function == "createPolicy" {
            return t.createPolicy(stub, args)
    }else if function == "calculateTokens" {
            return t.calculateTokens(stub, args)
    }
    fmt.Println("invoke did not find func: " + function)
    errMsg:= "{ \"message\" : \"Received unknown function invocation\", \"code\" : \"503\"}"
    err:= stub.SetEvent("errEvent", [] byte(errMsg))
    if err != nil {
        return nil, err
    }
    return nil, nil //error
}

// Query - Our entry Dealint for Queries
// ============================================================================================================================
func(t * Insurance) Query(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {
    fmt.Println("query is running " + function)
    // Handle different functions
    if function == "getPolicyByNumber" {
            return t.getPolicyByNumber(stub, args)
    }else if function == "getPolicyByPremium" {
            return t.getPolicyByPremium(stub, args)
    }else if function == "getAllPolicies" {
            return t.getAllPolicies(stub, args)
    }
    fmt.Println("query did not find func: " + function) //errors
    errMsg:= "{ \"message\" : \"Received unknown function query\", \"code\" : \"503\"}"
    err:= stub.SetEvent("errEvent", [] byte(errMsg))
    if err != nil {
        return nil, err
    }
    return nil, nil
}

// createPolicy stores the Policy on the ledger. If the policyNumber exists,
// it will return an error message
func(t * Insurance) createPolicy(stub shim.ChaincodeStubInterface, args []string) ([] byte, error) {

    if len(args) != 8 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 8\", \"code\" : \"503\"}"
        err := stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }

    fmt.Println("Create Policy Starting...")

    policyNumber := args[0]
    vehicleNumber := args[1]
    carUser := args[2]
    premiumNumber := args[3]
    carYear := args[4] 
    carModel := args[5]
    insuranceToken := args[6]
    status := args[7] 

    policyAsBytes, err:= stub.GetState(policyNumber)
    if err != nil {
        return nil, err
    }
    res:= Policy{}
    json.Unmarshal(policyAsBytes, &res)

    if res.PolicyNumber == policyNumber {
        fmt.Println("This Policy Number already exists: " + policyNumber)
        errMsg:= "{ \"policyNumber\" : \"" + policyNumber + "\", \"message\" : \"This policyNumber already exists\", \"code\" : \"503\"}"
        err:= stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil //all stop a Deal by this name exists
    }
    //build the Policy json string manually
    policyDetails := `{` + 
        `"policyNumber": "` + policyNumber + `" , ` + 
        `"vehicleNumber": "` + vehicleNumber + `" , ` + 
        `"carUser": "` + carUser + `" , ` + 
        `"premiumNumber": "` + premiumNumber + `" , ` + 
        `"carYear": "` + carYear + `" , ` + 
        `"carModel": "` + carModel + `" , ` + 
        `"insuranceToken": "` + insuranceToken + `" , ` + 
        `"status": "` + status + `" ` + 
    `}`
    fmt.Println("policyDetails: ");
    fmt.Println(policyDetails);
    err = stub.PutState(policyNumber, [] byte(policyDetails)) //store policy with policyNumber as key
    if err != nil {
        return nil, err
    }
    //get the policy index string
    policyIndexAsBytes, err:= stub.GetState(policyIndexStr)
    if err != nil {
        return nil, errors.New("Failed to get Policy index string")
    }

    fmt.Print("policyIndexAsBytes: ")
    fmt.Println(policyIndexAsBytes)
    var policyIndex[] string
    json.Unmarshal(policyIndexAsBytes, &policyIndex) //un stringify it aka JSON.parse()
    fmt.Print("policyIndex after unmarshal..before append: ")
    fmt.Println(policyIndex)
    //append
    policyIndex = append(policyIndex, policyNumber) //add policy policyNumber to index list
    fmt.Println("! policy index after appending policyNumber: ", policyIndex)
    jsonAsBytes, _:= json.Marshal(policyIndex)
    fmt.Print("jsonAsBytes: ")
    fmt.Println(jsonAsBytes)
    err = stub.PutState(policyIndexStr, jsonAsBytes) //store name of policy
    if err != nil {
        return nil, err
    }
    tosend:= "{ \"policyNumber\" : \"" + policyNumber + "\", \"message\" : \"policy created succcessfully\", \"code\" : \"200\"}"
        err = stub.SetEvent("evtsender", [] byte(tosend))
        if err != nil {
            return nil, err
        }
    return nil, nil
}

// getPolicyByNumber returns the value of the specified policyNumber
func(t * Insurance) getPolicyByNumber(stub shim.ChaincodeStubInterface, args []string) ([] byte, error) {
    if len(args) != 1 {
        errMsg:= "{ \"message\" : \"Incorrect arguments. Expecting a \"policyNumber\" as an argument\", \"code\" : \"503\"}"
        err := stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // set policyNumber
    policyNumber := args[0]

    value, err := stub.GetState(policyNumber)
    if err != nil {
        errMsg:= "{ \"message\" : \"" + policyNumber + " not Found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    
    return value, nil
}

func(t * Insurance) getPolicyByPremium(stub shim.ChaincodeStubInterface, args []string) ([] byte, error) {
    if len(args) != 1 {
        errMsg:= "{ \"message\" : \"Incorrect arguments. Expecting a \"premiumNumber\" as an argument\", \"code\" : \"503\"}"
        err := stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // set premiumNumber
    premiumNumber := args[0]

    //get the policy index string
    policyIndexAsBytes, err:= stub.GetState(policyIndexStr)
    if err != nil {
        return nil, err
    }

    fmt.Print("policyIndexAsBytes: ")
    fmt.Println(policyIndexAsBytes)
    var policyIndex[] string
    json.Unmarshal(policyIndexAsBytes, &policyIndex) //un stringify it aka JSON.parse()
    fmt.Println("policyIndex: ")
    fmt.Println(policyIndex)
    var valIndex Policy
    jsonResp := "{"
    for i, val:= range policyIndex {
        fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getPolicyByPremium")
        valueAsBytes, err:= stub.GetState(val)
        if err != nil {
            return nil, err
        }
        fmt.Print("valueAsBytes : ")
        fmt.Println(valueAsBytes)
        json.Unmarshal(valueAsBytes, &valIndex)
        fmt.Print("valIndex: ")
        fmt.Print(valIndex)
        if valIndex.PremiumNumber == premiumNumber {
            fmt.Println("PremiumNumber found: " + val)
            jsonResp = jsonResp + "\"" + val + "\":" + string(valueAsBytes[: ])
            fmt.Println("jsonResp inside if")
            fmt.Println(jsonResp)
            if i < len(policyIndex) - 1 {
                jsonResp = jsonResp + ","
            }
        } 
    }
    jsonResp = jsonResp + "}"
    fmt.Println("jsonResp : " + jsonResp)
    if jsonResp == "{}" {
        fmt.Println("PremiumNumber not found.")
    }
    
    return []byte(jsonResp), nil
}

func(t * Insurance) getAllPolicies(stub shim.ChaincodeStubInterface, args []string) ([] byte, error) {
    var jsonResp string
    var policyIndex[] string
    fmt.Println("getting all policies")
    var err error
    if len(args) != 1 {
        errMsg:= "{ \"message\" : \"Incorrect arguments. Expecting a \" \" as an argument\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    policyAsBytes, err:= stub.GetState(policyIndexStr)
    if err != nil {
        return nil, err
    }
    fmt.Print("policyAsBytes : ")
    fmt.Println(policyAsBytes)
    json.Unmarshal(policyAsBytes, &policyIndex) //un stringify it aka JSON.parse()
    fmt.Print("policyIndex : ")
    fmt.Println(policyIndex)

    jsonResp = "{"
    for i, val:= range policyIndex {
        fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all policies")
        valueAsBytes, err:= stub.GetState(val)
        if err != nil {
            return nil, err
        }
        fmt.Print("valueAsBytes : ")
        fmt.Println(valueAsBytes)
        jsonResp = jsonResp + "\"" + val + "\":" + string(valueAsBytes[: ])
        if i < len(policyIndex) - 1 {
            jsonResp = jsonResp + ","
        }
    }
    jsonResp = jsonResp + "}"
    fmt.Println("jsonResp : " + jsonResp)
    fmt.Println("Fetched all policies.")
    return []byte(jsonResp), nil
}

func(t * Insurance) calculateTokens(stub shim.ChaincodeStubInterface, args []string) ([] byte, error) {
    
    policyIndex:= Policy{}
    if len(args) != 4 {
        errMsg:= "{ \"message\" : \"Incorrect arguments. Expecting 4 arguments\", \"code\" : \"503\"}"
        err := stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }

    distance:= args[0]
    time:=args[1]
    numOfHours := args[2]
    policyNumber:= args[3]

    deductPercentage:= 0.00

    policyAsBytes, err:= stub.GetState(policyNumber)
    if err != nil {
        return nil, err
    }
    fmt.Print("policyAsBytes : ")
    fmt.Println(policyAsBytes)
    json.Unmarshal(policyAsBytes, &policyIndex) //un stringify it aka JSON.parse()
    fmt.Print("policyIndex : ")
    fmt.Println(policyIndex)
    _distance,err := strconv.Atoi(distance)
    if err != nil {
        fmt.Sprintf("Error while converting string 'distance' to int : %s", err.Error())
        return nil, errors.New("Error while converting string 'distance' to int ")
    }
    _time,err := strconv.Atoi(time)
    if err != nil {
        fmt.Sprintf("Error while converting string 'time' to int : %s", err.Error())
        return nil, errors.New("Error while converting string 'time' to int ")
    }
    _numOfHours,err := strconv.Atoi(numOfHours)
    if err != nil {
        fmt.Sprintf("Error while converting string 'numOfHours' to int : %s", err.Error())
        return nil, errors.New("Error while converting string 'numOfHours' to int ")
    }
    if _distance > 80{
        deductPercentage = 0.10   
    }
    if _time > 24 {
        deductPercentage = deductPercentage + 0.30
    }
    if _numOfHours >= 8{
        deductPercentage = deductPercentage + 0.40
    }
    fmt.Printf("deductPercentage: %s", deductPercentage)

    availableToken, errBool := strconv.ParseFloat(policyIndex.InsuranceToken, 64)
    if errBool != nil {
        fmt.Println(errBool)
    }

    tokenLeft := availableToken - deductPercentage
    policyIndex.InsuranceToken= strconv.FormatFloat(tokenLeft, 'f', 2, 64)

    policyDetails := `{` + 
        `"policyNumber": "` + policyIndex.PolicyNumber + `" , ` + 
        `"vehicleNumber": "` + policyIndex.VehicleNumber + `" , ` + 
        `"carUser": "` + policyIndex.CarUser + `" , ` + 
        `"premiumNumber": "` + policyIndex.PremiumNumber + `" , ` + 
        `"carYear": "` + policyIndex.CarYear + `" , ` + 
        `"carModel": "` + policyIndex.CarModel + `" , ` + 
        `"insuranceToken": "` + policyIndex.InsuranceToken + `" , ` + 
        `"status": "` + policyIndex.Status + `" ` + 
    `}`

    fmt.Println(policyDetails);
    err = stub.PutState(policyNumber, [] byte(policyDetails)) //store policy with policyNumber as key
    if err != nil {
        return nil, err
    }
    
    return []byte(policyIndex.InsuranceToken),nil
}
// main function starts up the chaincode in the container during instantiate
func main() {
    err := shim.Start(new(Insurance))
    if err != nil {
            fmt.Printf("Error starting Insurance chaincode: %s", err)
    }
}