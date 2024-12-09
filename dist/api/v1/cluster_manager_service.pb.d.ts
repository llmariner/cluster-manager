import * as fm from "../../fetch.pb";
import * as GoogleProtobufEmpty from "../../google/protobuf/empty.pb";
export type Cluster = {
    id?: string;
    name?: string;
    registrationKey?: string;
    object?: string;
    componentStatuses?: {
        [key: string]: ComponentStatus;
    };
};
export type CreateClusterRequest = {
    name?: string;
};
export type ListClustersRequest = {};
export type ListClustersResponse = {
    object?: string;
    data?: Cluster[];
};
export type GetClusterRequest = {
    id?: string;
};
export type DeleteClusterRequest = {
    id?: string;
};
export type DeleteClusterResponse = {
    id?: string;
    object?: string;
    deleted?: boolean;
};
export type InternalCluster = {
    cluster?: Cluster;
    tenantId?: string;
};
export type ListInternalClustersRequest = {};
export type ListInternalClustersResponse = {
    clusters?: InternalCluster[];
};
export type GetSelfClusterRequest = {};
export type ComponentStatus = {
    isHealthy?: boolean;
    message?: string;
};
export type UpdateComponentStatusRequest = {
    name?: string;
    status?: ComponentStatus;
};
export declare class ClustersService {
    static CreateCluster(req: CreateClusterRequest, initReq?: fm.InitReq): Promise<Cluster>;
    static ListClusters(req: ListClustersRequest, initReq?: fm.InitReq): Promise<ListClustersResponse>;
    static GetCluster(req: GetClusterRequest, initReq?: fm.InitReq): Promise<Cluster>;
    static DeleteCluster(req: DeleteClusterRequest, initReq?: fm.InitReq): Promise<DeleteClusterResponse>;
}
export declare class ClustersInternalService {
    static ListInternalClusters(req: ListInternalClustersRequest, initReq?: fm.InitReq): Promise<ListInternalClustersResponse>;
}
export declare class ClustersWorkerService {
    static GetSelfCluster(req: GetSelfClusterRequest, initReq?: fm.InitReq): Promise<Cluster>;
    static UpdateComponentStatus(req: UpdateComponentStatusRequest, initReq?: fm.InitReq): Promise<GoogleProtobufEmpty.Empty>;
}
