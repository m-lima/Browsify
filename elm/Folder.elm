import Http
import Html exposing (Html, button, div, text)
import Html.Events exposing (onClick)
import Json.Decode as Json exposing (..)

main = Html.beginnerProgram { model = model, view = view, update = update }

-- Model

type alias User =
  { name: String
  , email: String
}

userDecoder =
  map2 User (field "name" string) (field "email" string)

type alias Entry =
  { name: String
  , isDir: Bool
  , size: Int
  , date: String
}

entryDecoder =
  map4 Entry (field "name" string) (field "isDir" bool) (field "size" int) (field "date" string)

type alias Model =
  { user: Maybe User
  , folder: List Entry
}

model : Model
model =
  Model Nothing []

-- Update

type Msg = Increment
         | Decrement
         | Fetch
         | Logout

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
    Fetch ->
      Model (Just (Http.get userDecoder "https://localhost/user")) []
    Logout ->
      Http.post "https://localhost/logout" () ()
      Model (Just (Http.get userDecoder "https://localhost/user")) []


-- View

view : Model -> Html Msg
view model =
  case model.user of
    Nothing ->
      div []
      [ button [ onClick Fetch ] [ text "Fetch User" ]
      , div [] [ text "" ]
      ]
    Just user ->
      div []
      [ button [ onClick Logout ] [ text "Logout" ]
      , div [] [ text user.name ]
      ]
  -- case model.user of
  --   Nothing ->
  --     div []
  --     [ button [ onClick Decrement ] [ text "-" ]
  --     , div [] [ text "" ]
  --     , button [ onClick Increment ] [ text "+" ]
  --     ]
  --   Just user ->
  --     div []
  --     [ button [ onClick Decrement ] [ text "-" ]
  --     , div [] [ text user.name ]
  --     , button [ onClick Increment ] [ text "+" ]
  --     ]

