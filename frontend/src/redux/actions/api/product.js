import { API_URL, Client } from "./api";

export function listProduct(searchId) {
  const client = Client(true);
  return async (dispatch) => {
    try {
      let res = await client.get(`${API_URL}/product/${searchId}`);
      dispatch({
        type: "FETCH_LIST_PRODUCT",
        products: res.data,
      });
    } catch (err) {
      console.log(err);
    }
  };
}

export function resetProducts() {
  return (dispatch) => {
      dispatch({
        type: "RESET_LIST_PRODUCT",
        products: [],
      });
  };
}