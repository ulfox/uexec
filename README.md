# uexec

Golang lib for handling function return codes

## Generic Example

Here if Exec catches an error it will automatically do an os.Exit(1)

Content of **config.json**

    cat config.json
    {
        "version": "0.0.1",
        "git_sha": "b0d9c131648f6dfde6f1d976a58fb4ad0f55029a",
        "config": {
            "unit": null,
            "integration": null,
            "security": null
        }
    }

Marshal version into a struct

    package main

    import (
        "encoding/json"
        "fmt"
        "io/ioutil"

        "github.com/ulfox/uexec"
    )

    type Config struct {
        Version string
    }

    func main() {
        state := Config{}

        try := uexec.NewErrorHandler()
        try.Exec(
            json.Unmarshal(
                try.Exec(ioutil.ReadFile("config.json")
            ).ByteS(0), &state)
        )
        fmt.Println(state)
    }

Output

    0 > go run main.go 
    1   {0.0.1}

## Action Callback

Action callback is a function that can be defined to run right after
the called **function/method**.

### Callback Definition

The action callback must have the following definition

    func(...interface{}) interface{}

it takes any number of arguments and returns an interface

### Exampe with Action Callback


    package main

    import (
        "fmt"

        "github.com/ulfox/uexec"
    )

    func someCallBackFunction(cmd ...interface{}) interface{} {
        // some code
        return cmd[1]
    }

    func someFunction() (string, int, int, error) {
        // some code
        return "a", 0, 1, nil
    }

    func main() {
        erH := uexec.NewErrorHandler()
        output := erH.Exec(
            someFunction(),
        ).AddCallBack(
            someCallBackFunction,
            "A",
            "CallBack",
            "Function",
        ).CallBack()

        fmt.Println(output.Get(2))
        fmt.Println(output)
    }


First we initiate a new error handler **erH** and we execute **someFunction**.
What follows after is the **AddCallBack** method which registers a function **someCallBackFunction** and
its arguments **["A", "CallBack", "Function"]**. Last, we call the function callback function

Output

    0 > go run main.go 
    1   1
    2   {[a 0 1 <nil>] 0x4d2260 [a 0 1 <nil> A CallBack Function] 0 <nil>}

The first line after we run the code is the 3 output from **someFunction**.
The second line has:
 - **output.Values**: [a 0 1 <nil>]
 - **output.CallBackFunc**: 0x4d2260
 - **output.CallBackArgs**: [a 0 1 <nil> A CallBack Function]
 - **output.CallBackValues**: 0
 - **output.Err**: <nil>


### Elasticity

Elasticity can be defined as **true** in order to inform Exec not to exit on errors.

In a case where you simply want to run an action and continue regardless the outcome, then
**.SetElasticity(true)** needs to be set.

    package main

    import (
        "errors"
        "fmt"

        "github.com/ulfox/uexec"
    )

    func someFunction() (string, int, int, error) {
        // some code
        return "a", 0, 1, errors.New("someError")
    }

    func main() {
        erH := uexec.NewErrorHandler().SetElasticity(true)
        output := erH.Exec(
            someFunction(),
        )

        fmt.Println(output.Err)
    }

Notice that above we changed the nil error in the **someFunction** to 
errors.New("someError").

Output:

    0 > go run main.go 
    1   {"level":"error","msg":"someError","time":"2020-10-30T03:40:13+02:00"}
    2   someError

### Error Point - Define your error manually

Let's say a function retruns 3 values, from which none is an error.


    package main

    import (
        "fmt"

        "github.com/ulfox/uexec"
    )

    func someFunction() (*string, *int, *int) {
        // some code
        var v1 string = "a"
        var v2 int = 0
        var v3 int = 1
        return &v1, &v2, &v3
    }

    func main() {
        erH := uexec.NewErrorHandler().SetElasticity(true).EnableReportCaller(true)
        output := erH.ErP(2).Exec(
            someFunction(),
        )

        fmt.Println(output.Err)
    }


Notice that above we are returning 3 pointers.

The above example will cause an error because we are using **.ErP(2)** which instructs **Exec**
to check the index 2 (3rd) value and log an error in case that is not **<nil>**

Output:

    0 > go run main.go 
    1   {"file":"/datafs/.../github.com/ulfox/uexec@v0.../uexec.go:230","func":"....(*ErrorHandler).checkE","level":"error","msg":"0xc0000c0048","time":"2020-10-30T03:53:52+02:00"}
    2   0xc000014140 <-- Pointer of error


If elasticity was false, the error would have caused the program to exit immediately insteads of 
continuing to print the error.

## Generic Callback - Error Handling

The Generic Callback is defined the same way Action Callback is but is meant to run during an error.

    package main

    import (
        "errors"
        "fmt"

        "github.com/ulfox/uexec"
    )

    func someErrorCallBackFunction(cmd ...interface{}) interface{} {
        fmt.Println(cmd...)
        return cmd
    }

    func someFunction() (string, int, int, error) {
        // some code
        var v1 string = "a"
        var v2 int = 0
        var v3 int = 1
        return v1, v2, v3, errors.New("someError")
    }

    func main() {
        erH := uexec.NewErrorHandler().AddGenericCallBack(
            someErrorCallBackFunction,
            "generic",
            "error",
            "callback",
        ).ErP(3)

        output := erH.OnErr("callback").Exec(
            someFunction(),
        )

        fmt.Println(output.Err)
    }


Above we defined a new function **someErrorCallBackFunction** that simply prints the arguments it was called with.

Then we call **.AddGenericCallBack** method with the first argument being the function and the rest its arguments.

Output:

    0 > go run main.go 
    1  [a 0 1 someError generic error callback]
    2  someError


Instead of an exit we got in the the output (line 1) of the total arguments that were given in the error callback function.
