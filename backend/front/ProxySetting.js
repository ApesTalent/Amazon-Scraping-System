import React from "react";
import { Button, Input,  Modal,ModalHeader, ModalBody, ModalFooter} from "reactstrap";
import { connect } from "react-redux";
import {
  createProxy,
  createProxies,
  listProxy,
  deleteProxy,
  deleteProxies,
} from "../../../redux/actions/api";
import { Link } from "react-router-dom";
import * as Icon from "react-feather";

class Proxy extends React.Component {
  state = {
    newProxy: "",
    newProxies: "",
    showProxyModal: false,
  };

  componentDidMount = () => {
    this.props.listProxy();
  };

  toggleModal = () => {
    console.log("toggleModal")
    this.setState({ showProxyModal: !this.state.showProxyModal });
  };


  onChangeProxy = (e) => {
    this.setState({ newProxy: e.target.value });
  };


  onChangeProxies = (e) => {
    this.setState({ newProxies: e.target.value });
  };


  deleteProxy = async (id) => {
    await this.props.deleteProxy(id);
    this.props.listProxy();
  };

  deleteProxies = async () => {
    await this.props.deleteProxies();
    this.props.listProxy();
  };

  onAddProxy = async () => {
    const { newProxy } = this.state;
    if (!newProxy) return;
    await this.props.createProxy(newProxy);
    this.props.listProxy();
    this.setState({ newProxy: "" });
  };


  onAddProxies = async () => {
    const { newProxies } = this.state;
    if (!newProxies) return;
    await this.props.createProxies(newProxies);
    this.props.listProxy();
    this.setState({ newProxies: "" });
    this.toggleModal();
  };


  render() {
    const { proxys } = this.props;
    const { showProxyModal } = this.state;
    return (
      <React.Fragment>
        <h1>Proxies</h1>
        <div className="proxy-list-block">
          {proxys.map((prx) => (
            <div className="proxy-item" key={prx.id}>
              <div>{prx.index}</div>
              <span className="proxy-link">{prx.proxy}</span>
              <Link to="#" title="delete" onClick={() => this.deleteProxy(prx.id)}>
                <Icon.Trash size={18} />
              </Link>
            </div>
          ))}
        </div>
        <div className="proxy-add-block">
          <Input
            onChange={this.onChangeProxy}
            value={this.state.newProxy}
            placeholder="http://127.0.0.1:5858"
          />
          <Button color="primary" onClick={this.onAddProxy}>
            Add Proxy
          </Button>
          <Button className="ml-2" color="primary" onClick={this.toggleModal}>
            Add Proxies
          </Button>
          <Button className="ml-2" color="danger" onClick={this.deleteProxies}>
            Remove All
          </Button>
        </div>
        <Modal isOpen={showProxyModal} toggle={this.toggleModal}>
          <ModalHeader toggle={this.toggleModal}>Add Proxies</ModalHeader>
          <ModalBody>
            <div style={{ fontSize: "1.2rem" }}>
            <Input
            onChange={this.onChangeProxies}
            value={this.state.newProxies}
            type="textarea"
            rows={10}
            placeholder="input proxy lists..."
          />
            </div>
          </ModalBody>
          <ModalFooter>
            <Button color="primary" onClick={this.onAddProxies}>
              Yes
            </Button>{" "}
            <Button color="secondary" onClick={this.toggleModal}>
              No
            </Button>
          </ModalFooter>
        </Modal>
      </React.Fragment>
    );
  }
}

function mapStateToProps(state) {
  return {
    proxys: state.proxy.proxys,
  };
}

export default connect(mapStateToProps, {
  createProxy,
  createProxies,
  listProxy,
  deleteProxy,
  deleteProxies,
})(Proxy);
