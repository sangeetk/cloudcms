package worker

// PingPeer pings tne another worker
func (w *Worker) PingPeer(preq *PingRequest) (*PingResponse, error) {
	return nil, nil
}

// SyncPeers syncs all other peers within the same cluster
func (w *Worker) SyncPeers(sreq *SyncRequest) (*SyncResponse, error) {
	return nil, nil
}
