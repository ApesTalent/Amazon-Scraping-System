import cookie from "react-cookies";

const token = cookie.load("token")

const INITIAL_STATE = {
  user: {},
  authenticated: !!token
};

export default function (state = INITIAL_STATE, action) {
  switch (action.type) {
    case "FETCH_CURRENT_USER":
      return {
        ...state,
        user: action.user,
        authenticated: true
      }
    case "LOG_OUT":
      return {
        ...state,
        user: {},
        authenticated: false
      }
    default:
      return state;
  }
}
