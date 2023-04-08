import React from "react";
import { Button, Form, FormGroup, Input, Label, Row, Col } from "reactstrap";
import { connect } from "react-redux";
import { getAppConfig, updateAppConfig } from "../../../redux/actions/api";

class General extends React.Component {
  state = {
    interval: 0,
    discount: 0,
    token: "",
    chat_id: "",
    formerror: "",
  };

  componentDidMount = async () => {
    const { getAppConfig } = this.props;
    await getAppConfig();
    this.setState(this.props.setting);
  };

  updateConfig = (e) => {
    e.preventDefault();
    const { interval, discount, token, chat_id } = this.state;
    if (!interval) {
      this.setState({ formerror: "Interval is required" });
      return;
    }
    if (!discount) {
      this.setState({ formerror: "Discount is required" });
      return;
    }
    this.setState({ formerror: "" });
    this.props.updateAppConfig({
      interval: Math.floor(interval),
      discount: parseFloat(discount),
      token,
      chat_id
    });
  };

  render() {
    const { interval, discount, token, chat_id, formerror } = this.state;
    return (
      <Form className="mt-2" onSubmit={this.updateConfig}>
        <Row>
          <Col sm="12">
            <FormGroup>
              <Label for="interval">Scraping Interval (second)</Label>
              <Input
                id="interval"
                type="number"
                value={interval}
                onChange={(e) => this.setState({ interval: e.target.value })}
              />
            </FormGroup>
          </Col>
          <Col sm="12">
            <FormGroup>
              <Label for="discount">Discount (percent)</Label>
              <Input
                id="discount"
                type="number"
                value={discount}
                onChange={(e) => this.setState({ discount: e.target.value })}
              />
            </FormGroup>
          </Col>
          <Col sm="12">
            <FormGroup>
              <Label for="telegram">Telegram ID</Label>
              <Input
                id="telegram"
                value={chat_id}
                onChange={(e) => this.setState({ chat_id: e.target.value })}
              />
            </FormGroup>
          </Col>
          <Col sm="12">
            <FormGroup>
              <Label for="btoken">Telegram Token</Label>
              <Input
                id="btoken"
                value={token}
                onChange={(e) => this.setState({ token: e.target.value })}
              />
            </FormGroup>
          </Col>
          <Col sm="12">
            <div className="form-error">{formerror}</div>
          </Col>
          <Col className="d-flex justify-content-start flex-wrap" sm="12">
            <Button.Ripple className="mr-50" type="submit" color="primary">
              Save Changes
            </Button.Ripple>
          </Col>
        </Row>
      </Form>
    );
  }
}

function mapStateToProps(state) {
  return {
    setting: state.setting.setting,
  };
}

export default connect(mapStateToProps, {
  getAppConfig,
  updateAppConfig,
})(General);
