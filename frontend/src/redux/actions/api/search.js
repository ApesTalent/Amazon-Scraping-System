import { API_URL, errorMessage, createNotification, Client } from "./api";
import { history } from "../../../history";

export function listSearch() {
  const client = Client(true);
  return async (dispatch) => {
    try {
      let res = await client.get(`${API_URL}/search`);
      dispatch({
        type: "FETCH_LIST_SEARCH",
        searches: res.data,
      });
    } catch (err) {
      console.log(err);
    }
  };
}

export function createSearch(link) {
  return async (dispatch) => {
    const client = Client(true);
    try {
      await client.post(`${API_URL}/search`, {
        url: link
      });
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}

export function getSearch(id) {
  return async (dispatch) => {
    const client = Client(true);
    try {
      const res = await client.get(`${API_URL}/search/${id}`);
      return res.data
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}

export function updateSearch(search) {
  return async (dispatch) => {
    const client = Client(true);
    try {
      await client.put(`${API_URL}/search`, search);
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}

export function deleteSearch(id) {
  return async (dispatch) => {
    const client = Client(true);
    try {
      await client.delete(`${API_URL}/search/${id}`);
      history.push("/")
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}