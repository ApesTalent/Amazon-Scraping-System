import React from "react";
import { connect } from "react-redux";
import moment from "moment";
import {
  Card,
  CardHeader,
  CardTitle,
  CardBody,
  Input,
  Row,
  Col,
  UncontrolledDropdown,
  DropdownMenu,
  DropdownItem,
  DropdownToggle,
  Button,
  Spinner,
} from "reactstrap";
import { AgGridReact } from "ag-grid-react";
import * as Icon from "react-feather";
import { Link } from "react-router-dom";
import { listSearch, createSearch } from "../../redux/actions/api";
import { history } from "../../history";
import "../../assets/scss/plugins/tables/_agGridStyleOverride.scss";

class Search extends React.Component {
  state = {
    newlink: "",
    loading: false,
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
        headerName: "Status",
        field: "status",
        filter: true,
        width: 150,
        minWidth: 120,
        cellRendererFramework: params => {
          return (
            <span>{this.showStatus(params.data.status)}</span>
          )
        }
      },
      {
        headerName: "Last Search",
        field: "last_search",
        filter: true,
        width: 200,
        minWidth: 180,
        cellRendererFramework: params => {
          return (
            <span>{params.data.last_search &&
              moment(params.data.last_search).format("YYYY-MM-DD HH:mm:ss")}</span>
          )
        }
      },
      {
        headerName: "Items",
        field: "items",
        filter: true,
        width: 150,
        minWidth: 120,
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

  componentDidMount = () => {
    this.listSearches();
  };

  listSearches = async () => {
    this.setState({ loading: true });
    await this.props.listSearch();
    this.setState({ loading: false });
  };

  onChangeLink = (e) => {
    this.setState({ newlink: e.target.value });
  };

  onAddNewLink = async () => {
    const { newlink } = this.state;
    if (!newlink) return;
    await this.props.createSearch(newlink);
    this.props.listSearch();
    this.setState({ newlink: "" });
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

  gotoSearchDetail = (id) => {
    history.push(`/search/${id}`);
  };

  onSelectRow = () => {
    var selectedRows = this.gridApi.getSelectedRows();
    if (selectedRows.length === 0) return
    history.push(`/search/${selectedRows[0].id}`);
  }

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
    const { searches } = this.props;
    const {
      loading,
      newlink,
      columnDefs,
      defaultColDef,
      pageSize,
    } = this.state;
    return (
      <React.Fragment>
        <Card>
          <CardHeader>
            <CardTitle>
              Search Links{" "}
              {!loading && (
                <Link to="#" onClick={this.listSearches} title="refresh">
                  <Icon.RotateCw size={16} className="fonticon-wrap" />
                </Link>
              )}
              {loading && (
                <Spinner
                  style={{ width: "1.2rem", height: "1.2rem" }}
                  color="primary"
                />
              )}
            </CardTitle>
          </CardHeader>
          <CardBody>
            <div className="ag-theme-material ag-grid-table">
              <div className="ag-grid-actions d-flex justify-content-between flex-wrap mb-1">
                <div className="sort-dropdown">
                  <UncontrolledDropdown className="ag-dropdown p-1">
                    <DropdownToggle tag="div">
                      1 - {pageSize} of {searches.length}
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
                rowData={searches}
                onGridReady={this.onGridReady}
                colResizeDefault={"shift"}
                animateRows={true}
                pagination={true}
                pivotPanelShow="always"
                paginationPageSize={pageSize}
                resizable={true}
                onSelectionChanged={this.onSelectRow}
              />
            </div>
          </CardBody>
        </Card>
        <Card className="mt-5">
          <CardHeader>
            <CardTitle>Add Search Link</CardTitle>
          </CardHeader>
          <CardBody>
            <Row>
              <Col md="6" sm="12" className="mt-1">
                <Input
                  onChange={this.onChangeLink}
                  value={newlink}
                  placeholder="Search Link"
                />
              </Col>
              <Col md="6" sm="12" className="mt-1">
                <Button.Ripple color="primary" onClick={this.onAddNewLink}>
                  Add
                </Button.Ripple>
              </Col>
            </Row>
          </CardBody>
        </Card>
      </React.Fragment>
    );
  }
}

function mapStateToProps(state) {
  return {
    searches: state.search.searches,
  };
}

export default connect(mapStateToProps, {
  listSearch,
  createSearch,
})(Search);
