import React from "react";
import { Button, Form, FormGroup, Input, Label, Row, Col } from "reactstrap";
import { connect } from "react-redux";
import { changePassword } from "../../../redux/actions/api";

class ChangePassword extends React.Component {
  state = {
    old_pass: "",
    new_pass: "",
    new_conf_pass: "",
    formerror: "",
  };

  onChangePassword = (e) => {
    e.preventDefault();
    const { old_pass, new_pass, new_conf_pass } = this.state;
    if (!old_pass) {
      this.setState({ formerror: "Old password is required" });
      return;
    }
    if (new_pass.length < 6) {
      this.setState({ formerror: "Password should be more than 6 in length" });
      return;
    }
    if (new_pass !== new_conf_pass) {
      this.setState({ formerror: "Password confirmation doesn't match" });
      return;
    }
    this.setState({ formerror: "" });
    this.props.changePassword({ old_pass, new_pass });
  };

  render() {
    const { old_pass, new_pass, new_conf_pass, formerror } = this.state;

    return (
      <React.Fragment>
        <Row className="pt-1">
          <Col sm="12">
            <Form onSubmit={this.onChangePassword}>
              <FormGroup>
                <Label for="old_pass">Old password</Label>
                <Input
                  id="old_pass"
                  type="password"
                  value={old_pass}
                  onChange={(e) => this.setState({ old_pass: e.target.value })}
                />
              </FormGroup>
              <FormGroup>
                <Label for="new_pass">New password</Label>
                <Input
                  id="new_pass"
                  type="password"
                  value={new_pass}
                  onChange={(e) => this.setState({ new_pass: e.target.value })}
                />
              </FormGroup>
              <FormGroup>
                <Label for="new_conf_pass">Confirm password</Label>
                <Input
                  id="new_conf_pass"
                  type="password"
                  value={new_conf_pass}
                  onChange={(e) =>
                    this.setState({ new_conf_pass: e.target.value })
                  }
                />
              </FormGroup>
              <div className="form-error">{formerror}</div>
              <div className="d-flex justify-content-start flex-wrap">
                <Button.Ripple
                  className="mr-1 mb-1"
                  color="primary"
                  type="submit"
                >
                  Save Changes
                </Button.Ripple>
                <Button.Ripple
                  className="mb-1"
                  color="danger"
                  type="reset"
                  outline
                >
                  Cancel
                </Button.Ripple>
              </div>
            </Form>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

function mapStateToProps(state) {
  return {
    user: state.auth.user,
  };
}

export default connect(mapStateToProps, {
  changePassword,
})(ChangePassword);
