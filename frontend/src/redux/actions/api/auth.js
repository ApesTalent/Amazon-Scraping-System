import { API_URL, errorMessage, createNotification, Client } from "./api";
import cookie from "react-cookies";
import { history } from "../../../history";

export function protectedUser() {
  const client = Client(true);
  return async (dispatch) => {
    try {
      let res = await client.get(`${API_URL}/protected`);
      dispatch({
        type: "FETCH_CURRENT_USER",
        user: res.data,
      });
    } catch (err) {
      cookie.remove("token", { path: "/" });
      dispatch({ type: "LOG_OUT" });
    }
  };
}

export function loginUser(user) {
  return async (dispatch) => {
    const client = Client();
    try {
      let res = await client.post(`${API_URL}/login`, user);
      const token = res.data;
      cookie.save("token", token, { path: "/" });

      const client2 = Client(true);
      res = await client2.get(`${API_URL}/protected`);
      dispatch({
        type: "FETCH_CURRENT_USER",
        user: res.data,
      });
      history.push("/");
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}

export function changePassword(payload) {
  return async (dispatch) => {
    const client = Client(true);
    try {
      await client.post(`${API_URL}/changepwd`, payload);
      createNotification("success", "Password is updated successfully");
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}

export function logOut() {
  return async (dispatch) => {
    cookie.remove("token", { path: "/" });
    dispatch({ type: "LOG_OUT" });
    // history.push("/login");
  };
}
