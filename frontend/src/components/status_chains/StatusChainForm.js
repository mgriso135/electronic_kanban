import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';

const StatusChainForm = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [statuses, setStatuses] = useState(null);
  const [isEdit, setIsEdit] = useState(false);
  const [availableStatuses, setAvailableStatuses] = useState([]);

  useEffect(() => {
    const fetchAvailableStatuses = async () => {
      try {
          const response = await api.get('/statuses');
          setAvailableStatuses(response.data);
      } catch(error) {
        console.error('Error fetching available statuses:', error);
      }
    }

    const fetchChainStatuses = async () => {
        if (id) {
            setIsEdit(true);
          try {
            const response = await api.get(`/status-chains/${id}`);
            setName(response.data.name);

          const statusesResponse = await api.get(`/status-chains/${id}/statuses`);
            setStatuses(statusesResponse.data);
          } catch (error) {
              console.error('Error fetching status chain and statuses:', error);
          }
        } else {
            setStatuses([]); // Initialize statuses to empty array for new status chain
        }
    };

    fetchAvailableStatuses();
    fetchChainStatuses(); // Fetch chain statuses also moved into useEffect
  }, [id]);



  const handleStatusChange = (e, statusId) => {
    const updatedStatuses = statuses.map(status => {
      if(status.status_id === statusId){
        return { ...status, customer_supplier: parseInt(e.target.value, 10) };
      }
      return status;
    });
    setStatuses(updatedStatuses);
  };

   const handleStatusOrderChange = (e, statusId) => {
    const updatedStatuses = statuses.map(status => {
      if(status.status_id === statusId){
        return { ...status, order: parseInt(e.target.value, 10) };
      }
      return status;
    });
      setStatuses(updatedStatuses)
  };

  const addStatus = (statusId) => {
    const selectedStatus = availableStatuses.find(status => status.status_id === statusId);
      if (selectedStatus && !statuses?.some(s => s.status_id === selectedStatus.status_id)) {
      setStatuses([...(statuses || []), {...selectedStatus, status_name: selectedStatus.name, status_color: selectedStatus.color, customer_supplier: 1, order: (statuses || []).length +1}]);
      }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

  const statusChainData = {
      status_chain: {
        name: name,
      },
    statuses: statuses?.map(status => ({
      status_id: status.status_id,
      order: status.order,
      customer_supplier: status.customer_supplier
    })) || [],
  };

  console.log("handleSubmit - isEdit:", isEdit); // DEBUG: Log isEdit flag
  console.log("handleSubmit - statusChainData Payload:", statusChainData); // DEBUG: Log the entire payload

  try {
    if (isEdit) {
        console.log("handleSubmit - PUT request to update status chain:", `/status-chains/${id}`); // DEBUG: Log PUT request URL
        console.log("handleSubmit - PUT status_chain data:", statusChainData.status_chain); // DEBUG: Log status_chain part of payload
        await api.put(`/status-chains/${id}`, statusChainData.status_chain);
        console.log("handleSubmit - PUT request to update statuses:", `/status-chains/${id}/statuses`); // DEBUG: Log PUT request URL for statuses
        console.log("handleSubmit - PUT statuses data:", statusChainData.statuses); // DEBUG: Log statuses part of payload
        await api.put(`/status-chains/${id}/statuses`, statusChainData.statuses);
    } else {
       const response = await api.post('/status-chains', statusChainData);
          if (response.data && response.data.status_chain_id){
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
                        {statuses.map(status => (
                         <tr key={status.status_id}>
                           <td>{status.status_name}</td>
                           <td><input type="number" value={status.order} onChange={e => handleStatusOrderChange(e,status.status_id)} /></td>
                             <td><select value={status.customer_supplier} onChange={(e) => handleStatusChange(e, status.status_id)}>
                            <option value="1">Supplier</option>
                            <option value="2">Customer</option>
                           </select></td>
                          <td><button type="button" onClick={() => setStatuses(statuses.filter(s => s.status_id !== status.status_id))}>Delete</button></td>
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
                    .filter(status => !statuses?.some(s => s.status_id === status.status_id)) // Filter out already added statuses
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