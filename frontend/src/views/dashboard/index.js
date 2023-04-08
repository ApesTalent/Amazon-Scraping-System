import React from "react";
import { Row, Col } from "reactstrap";
import Search from "./Searches";
import RunStatus from "./RunStatus"

class Home extends React.Component {
  render() {
    return (
      <Row>
        <Col sm="12">
          <RunStatus />
        </Col>
        <Col sm="12">
          <Search />
        </Col>
      </Row>
    );
  }
}

export default Home;
