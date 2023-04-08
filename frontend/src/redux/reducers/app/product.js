const INITIAL_STATE = {
  products: [],
};

export default function (state = INITIAL_STATE, action) {
  switch (action.type) {
    case "FETCH_LIST_PRODUCT":
      let products = action.products || [];
      for (let i = 0; i < products.length; i++) {
        products[i].index = i + 1;
      }
      return {
        ...state,
        products: products,
      };
    default:
      return state;
  }
}
