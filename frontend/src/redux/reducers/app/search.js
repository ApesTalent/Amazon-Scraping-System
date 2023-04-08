const INITIAL_STATE = {
  searches: [],
};

export default function (state = INITIAL_STATE, action) {
  switch (action.type) {
    case "FETCH_LIST_SEARCH":
      let searches = action.searches || [];
      for (let i = 0; i < searches.length; i++) {
        searches[i].index = i + 1;
      }
      return {
        ...state,
        searches: searches,
      };
    default:
      return state;
  }
}
