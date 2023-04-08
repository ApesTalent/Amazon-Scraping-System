import React from "react";
import { connect } from "react-redux";
import {
  getSearch,
  listProduct,
  updateSearch,
  deleteSearch,
} from "../../redux/actions/api";
import moment from "moment";
import {
  Card,
  CardHeader,
  CardTitle,
  CardBody,
  Row,
  Col,
  UncontrolledDropdown,
  DropdownMenu,
  DropdownItem,
  DropdownToggle,
  Media,
  Button,
  Spinner,
  Modal,
  ModalHeader,
  ModalBody,
  ModalFooter,
} from "reactstrap";
import { AgGridReact } from "ag-grid-react";
import * as Icon from "react-feather";
import { Link } from "react-router-dom";
import "../../assets/scss/plugins/tables/_agGridStyleOverride.scss";

class Search extends React.Component {
  state = {
    search: {},
    loading: false,
    loadingScrape: false,
    showDeleteModal: false,
    pageSize: 5,
    isVisible: true,
    reload: false,
    defaultColDef: {
      sortable: true,
    },
    searchVal: "",
    columnDefs: [
      {
        headerName: "ID",
        field: "index",
        width: 100,
        suppressSizeToFit: true,
      },
      {
        headerName: "ASIN",
        field: "asin",
        width: 160,
        minWidth: 120
      },
      {
        headerName: "Title",
        field: "title",
        filter: true,
        width: 350,
        minWidth: 250,
        cellRendererFramework: params => {
          return (
            <div className="search-link" title={params.data.title}>{params.data.title}</div>
          )
        }
      },
      {
        headerName: "Price",
        field: "price",
        filter: true,
        width: 160,
        minWidth: 120
      },
      {
        headerName: "Link",
        field: "url",
        filter: true,
        width: 350,
        minWidth: 250,
        cellRendererFramework: params => {
          return (
            <div className="search-link" title={params.data.url}>{params.data.url}</div>
          )
        }
      },
    ],
  };

  componentDidMount = async () => {
    const id = this.props.match.params.searchId;
    const res = await this.props.getSearch(id);
    this.setState({ search: res });
    this.getProductList(id)
  };

  getProductList = async (searchID) => {
    this.setState({loading: true})
    await this.props.listProduct(searchID);
    this.setState({loading: false})
  }

  refreshProducts = () => {
    this.getProductList(this.state.search.id)
  };

  showStatus = (status) => {
    switch (status) {
      case 0:
        return "disable";
      case 1:
        return "active";
      case 2:
        return "searching";
      default:
        return "";
    }
  };

  toggleModal = () => {
    this.setState({ showDeleteModal: !this.state.showDeleteModal });
  };

  onUpdateSearch = async () => {
    const { search } = this.state;
    const { updateSearch } = this.props;
    const status = search.status === 1 ? 0 : 1;
    this.setState({ loadingScrape: true });
    search.status = status;
    await updateSearch(search);
    this.setState({ loadingScrape: false, search });
  };

  onDeleteSearch = async () => {
    const { search } = this.state;
    const { deleteSearch } = this.props;
    await deleteSearch(search.id);
    this.toggleModal();
  };

  onGridReady = (params) => {
    this.gridApi = params.api;
    this.gridColumnApi = params.columnApi;
    this.gridApi.sizeColumnsToFit()
  };

  filterData = (column, val) => {
    var filter = this.gridApi.getFilterInstance(column);
    var modelObj = null;
    if (val !== "all") {
      modelObj = {
        type: "equals",
        filter: val,
      };
    }
    filter.setModel(modelObj);
    this.gridApi.onFilterChanged();
  };

  filterSize = (val) => {
    if (this.gridApi) {
      this.gridApi.paginationSetPageSize(Number(val));
      this.setState({
        pageSize: val,
      });
    }
  };
  updateSearchQuery = (val) => {
    this.gridApi.setQuickFilter(val);
    this.setState({
      searchVal: val,
    });
  };

  render() {
    const { products } = this.props;
    const { search, loading, loadingScrape, showDeleteModal, defaultColDef, pageSize, columnDefs } = this.state;
    return (
      <React.Fragment>
        <Card>
          <CardHeader>
            <CardTitle>Search Item</CardTitle>
          </CardHeader>
          <CardBody>
            <Media body>
              <Row>
                <Col sm="9" md="6">
                  <div className="users-page-view-table">
                    <div className="d-flex user-info">
                      <div className="user-info-title font-weight-bold">
                        Status
                      </div>
                      <div>{this.showStatus(search.status)}</div>
                    </div>
                  </div>
                </Col>
                <Col sm="9" md="6">
                  <div className="users-page-view-table">
                    <div className="d-flex user-info">
                      <div className="user-info-title font-weight-bold">
                        Last Search
                      </div>
                      <div>
                        {search.last_search &&
                          moment(search.last_search).format(
                            "YYYY-MM-DD HH:mm:ss"
                          )}
                      </div>
                    </div>
                  </div>
                </Col>
              </Row>
              <Row>
                <Col>
                  <div className="users-page-view-table">
                    <div className="d-flex user-info">
                      <div className="user-info-title font-weight-bold">
                        Link
                      </div>
                      <div style={{ wordBreak: "break-all" }}>{search.url}</div>
                    </div>
                  </div>
                </Col>
              </Row>
              <Row className="mt-3">
                <Col className="search-btn-group">
                  <Button.Ripple
                    color="primary"
                    onClick={this.onUpdateSearch}
                  >
                    {loadingScrape && (
                      <Spinner
                        style={{ width: "1.2rem", height: "1.2rem" }}
                        color="light"
                      />
                    )}
                    {!loadingScrape && search.status === 1 && "Disable search"}
                    {!loadingScrape && search.status === 0 && "Activate search"}
                  </Button.Ripple>
                  <Button.Ripple
                    color="danger"
                    onClick={this.toggleModal}
                  >
                    Delete search
                  </Button.Ripple>
                </Col>
              </Row>
            </Media>
          </CardBody>
        </Card>
        <Card className="mt-5">
          <CardHeader>
            <CardTitle>
              Searched Products{" "}
              <Link to="#" onClick={this.refreshProducts} title="refresh">
                {!loading && <Icon.RotateCw size={16} className="fonticon-wrap" />}
                {loading && <Spinner style={{ width: '1.5rem', height: '1.5rem' }} color="primary" />}
              </Link>
            </CardTitle>
          </CardHeader>
          <CardBody>
            <div className="ag-theme-material ag-grid-table">
              <div className="ag-grid-actions d-flex justify-content-between flex-wrap mb-1">
                <div className="sort-dropdown">
                  <UncontrolledDropdown className="ag-dropdown p-1">
                    <DropdownToggle tag="div">
                      1 - {pageSize} of {products.length}
                      <Icon.ChevronDown className="ml-50" size={15} />
                    </DropdownToggle>
                    <DropdownMenu right>
                      <DropdownItem
                        tag="div"
                        onClick={() => this.filterSize(5)}
                      >
                        5
                      </DropdownItem>
                      <DropdownItem
                        tag="div"
                        onClick={() => this.filterSize(10)}
                      >
                        10
                      </DropdownItem>
                      <DropdownItem
                        tag="div"
                        onClick={() => this.filterSize(20)}
                      >
                        20
                      </DropdownItem>
                      <DropdownItem
                        tag="div"
                        onClick={() => this.filterSize(50)}
                      >
                        50
                      </DropdownItem>
                    </DropdownMenu>
                  </UncontrolledDropdown>
                </div>
              </div>
              <AgGridReact
                gridOptions={{}}
                rowSelection="single"
                defaultColDef={defaultColDef}
                columnDefs={columnDefs}
                rowData={products}
                onGridReady={this.onGridReady}
                colResizeDefault={"shift"}
                animateRows={true}
                pagination={true}
                pivotPanelShow="always"
                paginationPageSize={pageSize}
                resizable={true}
              />
            </div>
          </CardBody>
        </Card>
        <Modal isOpen={showDeleteModal} toggle={this.toggleModal}>
          <ModalHeader toggle={this.toggleModal}>Delete Search</ModalHeader>
          <ModalBody>
            <div className="p-3" style={{ fontSize: "1.2rem" }}>
              Do you really want to delete this search item?
            </div>
          </ModalBody>
          <ModalFooter>
            <Button color="primary" onClick={this.onDeleteSearch}>
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
    products: state.product.products,
  };
}

export default connect(mapStateToProps, {
  getSearch,
  listProduct,
  updateSearch,
  deleteSearch,
})(Search);
