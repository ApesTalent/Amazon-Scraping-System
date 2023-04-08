import { combineReducers } from "redux";
import customizer from "./customizer/";
import navbar from "./navbar/Index";
import search from "./app/search";
import product from "./app/product";
import auth from "./app/auth";
import setting from "./app/setting";
import proxy from "./app/proxy";

const rootReducer = combineReducers({
  customizer: customizer,
  auth: auth,
  navbar: navbar,
  search: search,
  product: product,
  setting: setting,
  proxy: proxy,
});

export default rootReducer;
