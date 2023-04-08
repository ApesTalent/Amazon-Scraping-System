const INITIAL_STATE = {
  proxys: [],
};

export default function (state = INITIAL_STATE, action) {
  switch (action.type) {
    case "FETCH_LIST_PROXY":
      let proxys = action.proxys || [];
      for (let i = 0; i < proxys.length; i++) {
        proxys[i].index = i + 1;
      }
      return {
        ...state,
        proxys: proxys,
      };
    default:
      return state;
  }
}
