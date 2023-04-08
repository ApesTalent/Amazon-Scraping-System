import { API_URL, errorMessage, createNotification, Client } from "./api";

export function listProxy() {
  const client = Client(true);
  return async (dispatch) => {
    try {
      let res = await client.get(`${API_URL}/proxy`);
      dispatch({
        type: "FETCH_LIST_PROXY",
        proxys: res.data,
      });
    } catch (err) {
      console.log(err);
    }
  };
}

export function createProxy(link) {
  return async (dispatch) => {
    const client = Client(true);
    try {
      await client.post(`${API_URL}/proxy`, {
        proxy: link
      });
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}


export function createProxies(links) {
  return async (dispatch) => {
    const client = Client(true);
    try {
      await client.post(`${API_URL}/proxies`, {
        proxy: links
      });
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}



export function deleteProxy(id) {
  return async (dispatch) => {
    const client = Client(true);
    try {
      await client.delete(`${API_URL}/proxy/${id}`);
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}

export function deleteProxies() {
  return async (dispatch) => {
    const client = Client(true);
    try {
      await client.delete(`${API_URL}/proxies`);
    } catch (err) {
      createNotification("error", errorMessage(err));
    }
  };
}