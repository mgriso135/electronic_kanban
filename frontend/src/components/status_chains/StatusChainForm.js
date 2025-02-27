import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';

const StatusChainForm = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [name, setName] = useState('');
    const [statuses, setStatuses] = useState([]);
    const [isEdit, setIsEdit] = useState(false);
    const [availableStatuses, setAvailableStatuses] = useState([]);


    useEffect(() => {
        const fetchAvailableStatuses = async () => {
            try {
                const response = await api.get('/statuses');
                setAvailableStatuses(response.data);
            } catch (error) {
                console.error('Error fetching available statuses:', error);
            }
        }

        fetchAvailableStatuses();

        if (id) {
            setIsEdit(true);
            const fetchStatusChain = async () => {
                try {
                    const response = await api.get(`/status-chains/${id}`);
                    setName(response.data.name);

                    const statusesResponse = await api.get(`/status-chains/${id}/statuses`);
                    setStatuses(statusesResponse.data);
                } catch (error) {
                    console.error('Error fetching status chain and statuses:', error);
                }
            };

            fetchStatusChain();
        } else {
            setStatuses([]); // Initialize statuses to empty array for new status chain
        }
    }, [id]);



    const handleStatusChange = (e, statusId) => {
        const updatedStatuses = statuses.map(status => {
            if (status.status_id === statusId) {
                return { ...status, customer_supplier: parseInt(e.target.value, 10) };
            }
            return status;
        });
        setStatuses(updatedStatuses);
    };

    const handleStatusOrderChange = (e, statusId) => {
        const updatedStatuses = statuses.map(status => {
            if (status.status_id === statusId) {
                return { ...status, order: parseInt(e.target.value, 10) };
            }
            return status;
        });
        setStatuses(updatedStatuses)
    };

    const addStatus = (statusId) => {
        const selectedStatus = availableStatuses.find(status => status.status_id === statusId);
        if (selectedStatus && !statuses?.some(s => s.status_id === selectedStatus.status_id)) {
            setStatuses([...(statuses || []), { ...selectedStatus, status_name: selectedStatus.name, customer_supplier: 1, order: (statuses || []).length + 1 }]);
        }
    };

    // New function to remove a status from the chain (mark for deletion)
    const removeStatusFromChain = (statusId) => {
        console.log("removeStatusFromChain CALLED - statusId:", statusId); // ADD LOG - Check if function is called and statusId value
        const updatedStatuses = statuses.map(status => {
            if (status.status_id === statusId) {
                return { ...status, _destroy: true }; // Mark for deletion using _destroy flag
            }
            return status;
        });
        setStatuses(updatedStatuses);
        console.log("removeStatusFromChain: setStatuses called with updatedStatuses:", updatedStatuses); // ADD LOG - Check updatedStatuses value
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        const statusesToDelete = statuses.filter(status => status._destroy); // Filter out statuses marked for deletion
        const statusesToUpdate = statuses.filter(status => !status._destroy); // Filter out statuses NOT marked for deletion

        const statusChainData = {
            status_chain: {
                name: name,
                ...(isEdit ? { status_chain_id: parseInt(id, 10) } : {})
            },
            statuses: statusesToUpdate.map(status => ({ // Send only statuses to update
                status_id: status.status_id,
                order: status.order,
                customer_supplier: status.customer_supplier
            })) || [],
        };

        console.log("StatusChainForm - handleSubmit: Sending statusChainData to API:", statusChainData);

        try {
            if (isEdit) {
                // **Send DELETE requests for statuses marked for deletion BEFORE PUT request:**
                for (const statusToDelete of statusesToDelete) {
                    console.log("StatusChainForm - handleSubmit: Deleting status ID:", statusToDelete.status_id, "from chain ID:", id);
                    await api.delete(`/status-chains/${id}/statuses/${statusToDelete.status_id}`); // New DELETE API call
                }


                await api.put(`/status-chains/${id}`, statusChainData.status_chain);
                await api.put(`/status-chains/${id}/statuses`, statusChainData.statuses); // Send updates for remaining statuses

            } else {
                const response = await api.post('/status-chains', statusChainData);
                if (response.data && response.data.status_chain_id) {
                    await api.put(`/status-chains/${response.data.status_chain_id}/statuses`, statusChainData.statuses);
                }
            }
            navigate('/status-chains');
        } catch (error) {
            console.error('Error saving status chain:', error);
        }
    };


    return (
        <div>
            <h2>{isEdit ? 'Edit Status Chain' : 'Create New Status Chain'}</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Name:</label>
                    <input type="text" value={name} onChange={(e) => setName(e.target.value)} />
                </div>

                <h3>Statuses in Chain</h3>

                {statuses != null && statuses.length > 0 ?
                    <table>
                        <thead>
                            <tr>
                                <th>Status Name</th>
                                <th>Order</th>
                                <th>Customer/Supplier</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                        {statuses
                         .filter(status => !status._destroy) // **ADD FILTER HERE - Only render non-deleted statuses**
                         .map(status => (
                         <tr key={status.status_id}>
                           <td>{status.status_name}</td>
                           <td><input type="number" value={status.order} onChange={e => handleStatusOrderChange(e,status.status_id)} /></td>
                             <td>
                            <select value={status.customer_supplier} onChange={(e) => handleStatusChange(e, status.status_id)}>
                            <option value="1">Supplier</option>
                            <option value="2">Customer</option>
                           </select>
                          </td>
                         <td>
                           <button type="button" onClick={() => removeStatusFromChain(status.status_id)}>
                              Delete
                           </button>
                         </td>
                        </tr>
                      ))}
                      </tbody>
                    </table>
                    :
                    <p>No Statuses Selected</p>

                }


                <div>
                    <h4>Available Statuses</h4>
                    <select onChange={(e) => addStatus(parseInt(e.target.value, 10))}>
                        <option value="">Select a Status</option>
                        {availableStatuses
                            .filter(status => !statuses?.some(s => s.status_id === status.status_id))
                            .map(status => (
                                <option key={status.status_id} value={status.status_id}>{status.name}</option>
                            ))}
                    </select>
                </div>
                <button type="submit">{isEdit ? 'Update Status Chain' : 'Create Status Chain'}</button>
                <Link to="/status-chains">Cancel</Link>
            </form>
        </div>
    );
};

export default StatusChainForm;