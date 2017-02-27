import Html exposing (Html, button, div, text)
import Html.Events exposing (onClick)

main = Html.beginnerProgram { model = model, view = view, update = update }

-- Model

type alias User =
  { name: String
  , email: String
}

type alias Entry =
  { name: String
  , isDir: Bool
  , size: Int
  , date: String
}

type alias Model =
  { user: Maybe User
  , folder: List Entry
}

model : Model
model =
  Model Nothing []

-- Update

type Msg = Increment | Decrement

update : Msg -> Model -> Model
update msg model =
  case msg of
    Increment ->
      case model.user of
        Nothing ->
          Model (Just (User "x" "")) []
        Just user ->
          Model (Just (User (user.name ++ "x") "")) []
    Decrement ->
      Model Nothing []

-- View

view : Model -> Html Msg
view model =
  case model.user of
    Nothing ->
      div []
      [ button [ onClick Decrement ] [ text "-" ]
      , div [] [ text "" ]
      , button [ onClick Increment ] [ text "+" ]
      ]
    Just user ->
      div []
      [ button [ onClick Decrement ] [ text "-" ]
      , div [] [ text user.name ]
      , button [ onClick Increment ] [ text "+" ]
      ]
