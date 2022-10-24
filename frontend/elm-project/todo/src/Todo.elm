module Todo exposing (main)
import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput)



--- MAIN
main = Browser.sandbox { init = init, update = update, view = view }
-- MODEL

type alias Model =
  { name : String
  , task : String
  }


init : Model
init =
  Model "" ""



-- UPDATE


type Msg
  = Name String
  | Task String

update : Msg -> Model -> Model
update msg model =
  case msg of
    Name name ->
      { model | name = name }

    Task task ->
      { model | task = task }

-- VIEW

view : Model -> Html Msg
view model =
  div []
        [ h1 [] [ text "Todo List" ]
        , Html.form []
            [ div []
        [ viewInput "text" "Name" model.name Name 
        , viewInput "text" "Task" model.task Task
        , viewValidation model
        ]
        , button []
            [ text "Submit" ]
        ]
        ]
  
viewInput : String -> String -> String -> (String -> msg) -> Html msg
viewInput t p v toMsg =
  input [ type_ t, placeholder p, value v, onInput toMsg ] []

viewValidation : Model -> Html msg
viewValidation model =
  if model.name == "" || model.task == "" then
    div [ style "color" "red", style "text-align" "center" ] [ text "Please Fill All Fields!" ]
  else
    div [ style "color" "green",  style "text-align" "center" ] [ text "Good!" ]