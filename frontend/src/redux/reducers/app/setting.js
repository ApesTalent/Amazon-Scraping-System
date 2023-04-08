const INITIAL_STATE = {
  setting: {},
};

export default function (state = INITIAL_STATE, action) {
  switch (action.type) {
    case "FETCH_APP_SETTING":
      return {
        ...state,
        setting: action.setting,
      };
    case "UPDATE_RUN_STATUS":
      let setting = state.setting;
      setting.run_scrape = action.status;
      return {
        ...state,
        setting,
      };
    default:
      return state;
  }
}
