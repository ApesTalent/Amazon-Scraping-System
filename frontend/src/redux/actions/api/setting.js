import { API_URL, errorMessage, createNotification, Client } from "./api";

export function getAppConfig() {
  const client = Client(true);
  return async (dispatch) => {
    try {
      let res = await client.get(`${API_URL}/setting`);
      dispatch({
        type: "FETCH_APP_SETTING",
        setting: res.data,
      });
    } catch (err) {
      console.log(err);
    }
  };
}

export function updateAppConfig(config) {
  return async (dispatch) => {
    const client = Client(true);
    try {
      await client.post(`${API_URL}/setting`, config);
      let res = await client.get(`${API_URL}/setting`);
      dispatch({
        type: "FETCH_APP_SETTING",
        setting: res.data,
      });
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}

export function updateRunStatus(status) {
  return async (dispatch) => {
    const client = Client(true);
    try {
      await client.put(`${API_URL}/setting/${status}`, {});
      dispatch({
        type: "UPDATE_RUN_STATUS",
        status,
      });
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}
