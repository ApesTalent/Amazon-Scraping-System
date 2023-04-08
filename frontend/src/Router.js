import React, { Suspense, lazy } from "react";
import { Router, Switch, Route } from "react-router-dom";
import { history } from "./history";
import { connect } from "react-redux";
import Spinner from "./components/@vuexy/spinner/Loading-spinner";
import { ContextLayout } from "./utility/context/Layout";
import { protectedUser, getAppConfig } from "./redux/actions/api";
import RequireAuth from "./views/pages/authentication/RequireAuth";

// Route-based code splitting
const Home = lazy(() => import("./views/dashboard"));
const Setting = lazy(() => import("./views/pages/account-settings/AccountSettings"));
const Search = lazy(() => import("./views/pages/Search"));
const login = lazy(() => import("./views/pages/authentication/Login"));
const NotFound = lazy(() => import("./views/pages/NotFound"));

// Set Layout and Component Using App Route
const RouteConfig = ({
  component: Component,
  fullLayout,
  permission,
  user,
  ...rest
}) => (
  <Route
    {...rest}
    render={(props) => {
      return (
        <ContextLayout.Consumer>
          {(context) => {
            let LayoutTag =
              fullLayout === true
                ? context.fullLayout
                : context.state.activeLayout === "horizontal"
                ? context.horizontalLayout
                : context.VerticalLayout;
            return (
              <LayoutTag {...props} permission={props.user}>
                <Suspense fallback={<Spinner />}>
                  <Component {...props} />
                </Suspense>
              </LayoutTag>
            );
          }}
        </ContextLayout.Consumer>
      );
    }}
  />
);

const AppRoute = connect(mapStateToProps)(RouteConfig);

class AppRouter extends React.Component {
  componentDidMount = async () => {
    await this.props.protectedUser()
    this.props.getAppConfig()
  }

  render() {
    return (
      // Set the directory path if you are deploying in sub-folder
      <Router history={history}>
        <Switch>
          <AppRoute exact path="/" component={RequireAuth(Home)} />
          <AppRoute exact path="/setting" component={RequireAuth(Setting)} />
          <AppRoute exact path="/search/:searchId" component={RequireAuth(Search)} />
          <AppRoute path="/login" component={login} fullLayout />
          <AppRoute exact path="/*" component={NotFound} />
        </Switch>
      </Router>
    );
  }
}

function mapStateToProps(state) {
  return {
    user: state.auth.user.role,
    auth: state.auth.authenticated,
  };
}

export default connect(mapStateToProps, {
  protectedUser,
  getAppConfig
})(AppRouter);
