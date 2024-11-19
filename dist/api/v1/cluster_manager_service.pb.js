/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
import * as fm from "../../fetch.pb";
export class ClustersService {
    static CreateCluster(req, initReq) {
        return fm.fetchReq(`/v1/clusters`, Object.assign(Object.assign({}, initReq), { method: "POST", body: JSON.stringify(req) }));
    }
    static ListClusters(req, initReq) {
        return fm.fetchReq(`/v1/clusters?${fm.renderURLSearchParams(req, [])}`, Object.assign(Object.assign({}, initReq), { method: "GET" }));
    }
    static GetCluster(req, initReq) {
        return fm.fetchReq(`/v1/clusters/${req["id"]}?${fm.renderURLSearchParams(req, ["id"])}`, Object.assign(Object.assign({}, initReq), { method: "GET" }));
    }
    static DeleteCluster(req, initReq) {
        return fm.fetchReq(`/v1/clusters/${req["id"]}`, Object.assign(Object.assign({}, initReq), { method: "DELETE" }));
    }
}
export class ClustersInternalService {
    static ListInternalClusters(req, initReq) {
        return fm.fetchReq(`/llmariner.clusters.server.v1.ClustersInternalService/ListInternalClusters`, Object.assign(Object.assign({}, initReq), { method: "POST", body: JSON.stringify(req) }));
    }
}
export class ClustersWorkerService {
    static GetSelfCluster(req, initReq) {
        return fm.fetchReq(`/llmariner.clusters.server.v1.ClustersWorkerService/GetSelfCluster`, Object.assign(Object.assign({}, initReq), { method: "POST", body: JSON.stringify(req) }));
    }
    static UpdateComponentStatus(req, initReq) {
        return fm.fetchReq(`/llmariner.clusters.server.v1.ClustersWorkerService/UpdateComponentStatus`, Object.assign(Object.assign({}, initReq), { method: "POST", body: JSON.stringify(req) }));
    }
}
