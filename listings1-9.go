/////////////////////////////////////
//Listing 1: Mit Funktionen spielen//
/////////////////////////////////////

// Deklaration.
type Converter func(string) string

// Fabrik.
func MakePrefixer(prefix string) Converter {
    return func(in string) string {
      return prefix + ":" + in
    }
}

type Repeater int

// Methode mit Signatur des Converter.
func (r Repeater) Repeat(in string) string {
    return strings.Repeat(in, int(r))
}

// Anwendung des Converters als Argument.
func ConvertAll(in []string, c Converter) []string {
    out := make([]string, len(in))
    for i, s := range in {
      out[i] = c(s)
    }
    return out
}

...

var r Repeater = 5

in := []string{"fun", "with", "functions"}

// Verschiedene Anwendungen des Converters.
foos := ConvertAll(in, MakePrefixer("foo"))
repeated := ConvertAll(in, r.Repeat)
upperCased := ConvertAll(in, strings.ToUpper)

-------

//////////////////////////////////////////
//Listing 2: 'Hello, World!' umständlich//
//////////////////////////////////////////

package main

import "fmt"

type Duck interface {
    Quack() string
}

type NamedDuck func() string

func (nd NamedDuck) Quack() string {
    return fmt.Sprintf("Quack, I'm %s.", nd())
}

func greet(duck Duck) {
    fmt.Println(duck.Quack())
    fmt.Println("Hello, and I'm the computer.")
}

func main() {
    var donald NamedDuck = func() string {
      return "Donald"
    }
    greet(donald)
}

-------

//////////////////////////////////////
//Listing 3: Flacher Programmverlauf//
//////////////////////////////////////

filename := "/path/to/my/file"
file, err := os.Open(filname)
if err != nil {
    return fmt.Errorf("cannot open %q: %v", filename, err)
}
defer file.Close()

// Weitermachen ...

-------

/////////////////////////////////////////
//Listing 4: Nebenläufige State Machine//
/////////////////////////////////////////

type Data struct { ... }
type State func(Data) (State, error)
type StateMachine struct {
    ctx  context.Context
    state State
    actions chan func()
}

func New(ctx context.Context) *StateMachine {
  sm := &StateMachine{
    ctx: ctx,
    actions: make(chan func()),
  }
  sm.state = sm.off
  go sm.loop()
  return sm
}

func (sm *StateMachine) loop() {
    for {
        select {
        case <-sm.ctx.Done():
          return
        case action := <-sm.actions:
          action()
        }
    }
}

func (sm *StateMachine) Process(d Data) error {
    var err error
    action := func() {
      var state State
      state, err = sm.state(d)
      if err == nil {
        sm.state = state
      }
    }
    select {
    case sm.actions <- action:
    case <-time.After(timeout):
      err = errors.New("timeout")
    }
    return err
}

-------

////////////////////////////////////////
//Listing 5: Einfache Handler-Funktion//
////////////////////////////////////////

var h http.HandlerFunc = func(
    w http.ResponseWriter, 
    r *http.Request,
) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hallo, liebe Leser!"))
}

http.ListenAndServe(":8080", h)

-------

/////////////////////////////////
//Listing 6: Multiplexer nutzen//
/////////////////////////////////

mux := http.NewServeMux()

mux.Handle("/network/", NewNetworkHandler())
mux.Handle("/network/addresses/", NewAddressesHandler())
mux.Handle("/network/interfaces/", NewInterfacesHandler())
mux.HandleFunc("/", func(
    w http.ResponseWriter, 
    r *http.Request,
) {
    w.WriteHeader(http.StatusNotFound)
})

http.ListenAndServe(":8080", mux)


-------

//////////////////////////////////////////////////////////////////
//Listing 7: Method Wrapper für die Verteilung nach HTTP-Methode//
//////////////////////////////////////////////////////////////////

type GetHandler interface {
    ServeGet(w http.ResponseWriter, r *http.Request)
}

type PostHandler interface {
    ServePost(w http.ResponseWriter, r *http.Request)
}
...
type MethodWrapper struct {
    handler http.Handler
}

func (mw MethodWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
      if h, ok := mw.handler.(GetHandler); ok {
        h.ServeGet(w, r)
        return
      }
    case http.MethodPost:
      if h, ok := mw.handler.(PostHandler); ok {
        h.ServePost(w, r)
        return
      }
    case ...:
      ...
    }
    mw.handler.ServeHTTP(w, r)
}


-------

/////////////////////////////////////////////////////
//Listing 8: Definition eines Structs mit JSON-Tags//
/////////////////////////////////////////////////////

type MeteringPointValue struct {
    Namespace string        `json:"namespace"`
    ID        string        `json:"id"`
    Quantity  int           `json:"quantity"`
    Total     time.Duration `json:"total-duration"`
    Minimum   time.Duration `json:"minimum-duration"`
    Maximum   time.Duration `json:"maximum-duration"`
    Average   time.Duration `json:"average-duration"`
}


-------

/////////////////////////////////////////////
//Listing 9: Daten im JSON-Format versenden//
/////////////////////////////////////////////

func (h MPVHandler) ServeHTTP(
    w http.ResponseWriter, 
    r *http.Request,
) {
    w.Header().Set(
      "Content-Type", 
      "application/json; charset=ISO-8859-1",
    )
    enc := json.NewEncoder(w)
    mpv := ...
    err := enc.Encode(mpv)
    if err != nil {
      http.Error(
        w, 
        err.Error(),
        http.StatusInternalServerError,
      )
    }
    w.WriteHeader(http.StatusOK)
}