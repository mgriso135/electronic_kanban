import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';

const KanbanChainForm = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [clienteId, setClienteId] = useState('');
  const [prodottoCodice, setProdottoCodice] = useState('');
  const [fornitoreId, setFornitoreId] = useState('');
  const [leadtimeDays, setLeadtimeDays] = useState('');
  const [quantity, setQuantity] = useState('');
  const [tipoContenitore, setTipoContenitore] = useState('Box');
    const [statusChainId, setStatusChainId] = useState('');
  const [noOfActiveKanbans, setNoOfActiveKanbans] = useState(1);
    const [isEdit, setIsEdit] = useState(false);
  const [availableStatusChains, setAvailableStatusChains] = useState([]);


  useEffect(() => {
    const fetchAvailableStatusChains = async () => {
      try {
        const response = await api.get('/status-chains');
        setAvailableStatusChains(response.data)
      } catch (error) {
         console.error("Error fetching available status chains:", error);
      }
    };
      fetchAvailableStatusChains();
    if (id) {
      setIsEdit(true);
        const fetchKanbanChain = async () => {
          try {
             const response = await api.get(`/kanban-chains/${id}`);
              setClienteId(response.data.cliente_id);
              setProdottoCodice(response.data.prodotto_codice);
              setFornitoreId(response.data.fornitore_id);
              setLeadtimeDays(response.data.leadtime_days);
              setQuantity(response.data.quantity);
              setTipoContenitore(response.data.tipo_contenitore);
              setStatusChainId(response.data.status_chain_id);
             setNoOfActiveKanbans(response.data.no_of_active_kanbans);
            } catch(error) {
              console.error('Error fetching kanban chain:', error);
            }
        };
        fetchKanbanChain();
    }
  }, [id]);



  const handleSubmit = async (e) => {
      e.preventDefault();
    const kanbanChainData = {
        kanban_chain: {
          cliente_id: parseInt(clienteId, 10),
          prodotto_codice: prodottoCodice,
          fornitore_id: parseInt(fornitoreId, 10),
          leadtime_days: parseInt(leadtimeDays, 10),
          quantity: parseFloat(quantity),
          tipo_contenitore: tipoContenitore,
          status_chain_id: parseInt(statusChainId, 10),
            no_of_active_kanbans: parseInt(noOfActiveKanbans, 10),
        },
      no_of_initial_kanbans: parseInt(noOfActiveKanbans, 10)
    };
    try {
      if (isEdit) {
          await api.put(`/kanban-chains/${id}`, kanbanChainData.kanban_chain);
      } else {
          await api.post('/kanban-chains', kanbanChainData);
      }
        navigate('/kanban-chains');
    } catch(error) {
      console.error('Error saving kanban chain:', error)
    }
  };

  return (
    <div>
      <h2>{isEdit ? 'Edit Kanban Chain' : 'Create New Kanban Chain'}</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Customer ID:</label>
          <input type="number" value={clienteId} onChange={(e) => setClienteId(e.target.value)} />
        </div>
        <div>
          <label>Supplier ID:</label>
          <input type="number" value={fornitoreId} onChange={(e) => setFornitoreId(e.target.value)} />
        </div>
        <div>
          <label>Product ID:</label>
          <input type="text" value={prodottoCodice} onChange={(e) => setProdottoCodice(e.target.value)} />
        </div>
          <div>
              <label>Lead Time (days):</label>
              <input type="number" value={leadtimeDays} onChange={(e) => setLeadtimeDays(e.target.value)} />
          </div>
          <div>
              <label>Container Type:</label>
              <select value={tipoContenitore} onChange={(e) => setTipoContenitore(e.target.value)}>
                  <option value="Box">Box</option>
                  <option value="Pallet">Pallet</option>
                  <option value="Other">Other</option>
              </select>
           </div>
          <div>
              <label>Quantity:</label>
              <input type="number" value={quantity} onChange={(e) => setQuantity(e.target.value)} />
          </div>
           <div>
              <label>Status Chain:</label>
              <select value={statusChainId} onChange={(e) => setStatusChainId(e.target.value)}>
                <option value="">Select Status Chain</option>
                {availableStatusChains.map(chain => (
                   <option key={chain.status_chain_id} value={chain.status_chain_id}>{chain.name}</option>
                ))}
              </select>
          </div>
          <div>
            <label>Number of Active Kanbans:</label>
            <input type="number" value={noOfActiveKanbans} onChange={(e) => setNoOfActiveKanbans(e.target.value)} />
          </div>

          <button type="submit">{isEdit ? 'Update Kanban Chain' : 'Create Kanban Chain'}</button>
          <Link to="/kanban-chains">Cancel</Link>
      </form>
    </div>
  );
};

export default KanbanChainForm;