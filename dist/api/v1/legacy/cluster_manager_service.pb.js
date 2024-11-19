/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
import * as fm from "../../../fetch.pb";
export class ClustersInternalService {
    static ListInternalClusters(req, initReq) {
        return fm.fetchReq(`/llmoperator.clusters.server.v1.ClustersInternalService/ListInternalClusters`, Object.assign(Object.assign({}, initReq), { method: "POST", body: JSON.stringify(req) }));
    }
}
export class ClustersWorkerService {
    static GetSelfCluster(req, initReq) {
        return fm.fetchReq(`/llmoperator.clusters.server.v1.ClustersWorkerService/GetSelfCluster`, Object.assign(Object.assign({}, initReq), { method: "POST", body: JSON.stringify(req) }));
    }
}
