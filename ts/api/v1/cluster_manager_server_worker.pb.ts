/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as GoogleProtobufEmpty from "../../google/protobuf/empty.pb"
import * as LlmarinerClustersServerV1Cluster_manager_server from "./cluster_manager_server.pb"
export type UpdateComponentStatusRequest = {
  name?: string
  status?: LlmarinerClustersServerV1Cluster_manager_server.ComponentStatus
}

export class ClustersWorkerService {
  static UpdateComponentStatus(req: UpdateComponentStatusRequest, initReq?: fm.InitReq): Promise<GoogleProtobufEmpty.Empty> {
    return fm.fetchReq<UpdateComponentStatusRequest, GoogleProtobufEmpty.Empty>(`/llmariner.clusters.server.v1.ClustersWorkerService/UpdateComponentStatus`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
}