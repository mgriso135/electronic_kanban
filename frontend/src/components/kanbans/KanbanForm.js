import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';

const KanbanForm = () => {
    const { id } = useParams(); // Kanban ID for edit mode
    const navigate = useNavigate();
    const [kanbanChainId, setKanbanChainId] = useState('');
    const [statusChainId, setStatusChainId] = useState('');
    const [statusCurrent, setStatusCurrent] = useState('');
    const [leadtimeDays, setLeadtimeDays] = useState('');
    const [tipoContenitore, setTipoContenitore] = useState('');
    const [quantity, setQuantity] = useState('');
    const [isActive, setIsActive] = useState(true);
    const [isEdit, setIsEdit] = useState(false);
    const [availableKanbanChains, setAvailableKanbanChains] = useState([]);
    const [availableStatuses, setAvailableStatuses] = useState([]);


    useEffect(() => {
        const fetchDropdownData = async () => {
            try {
                const kanbanChainsResponse = await api.get('/kanban-chains');
                setAvailableKanbanChains(kanbanChainsResponse.data);
                const statusesResponse = await api.get('/statuses');
                setAvailableStatuses(statusesResponse.data);

            } catch (error) {
                console.error("Error fetching dropdown data:", error);
            }
        };
        fetchDropdownData();

        if (id) {
            setIsEdit(true);
            const fetchKanban = async () => {
                try {
                    const response = await api.get(`/kanbans/${id}`);
                    const k = response.data;
                    setKanbanChainId(String(k.kanban_chain_id));
                    setStatusChainId(String(k.status_chain_id));
                    setStatusCurrent(String(k.status_current));
                    setLeadtimeDays(String(k.leadtime_days));
                    setTipoContenitore(k.tipo_contenitore);
                    setQuantity(String(k.quantity));
                    setIsActive(k.is_active);
                } catch (error) {
                    console.error('Error fetching kanban:', error);
                }
            };
            fetchKanban();
        }
    }, [id]);

    useEffect(() => {
        // Fetch Kanban Chain details when kanbanChainId changes
        if (kanbanChainId) {
            const fetchKanbanChainDetails = async () => {
                try {
                    const response = await api.get(`/kanban-chains/${kanbanChainId}`);
                    const kc = response.data;
                    setStatusChainId(String(kc.status_chain_id));
                    setLeadtimeDays(String(kc.leadtime_days));
                    setTipoContenitore(kc.tipo_contenitore);
                    setQuantity(String(kc.quantity));
                } catch (error) {
                    console.error('Error fetching kanban chain details:', error);
                    setStatusChainId(''); // Clear status chain if error
                    setLeadtimeDays('');
                    setTipoContenitore('');
                    setQuantity('');
                }
            };
            fetchKanbanChainDetails();
        } else {
            // Clear derived fields if no kanban chain is selected
            setStatusChainId('');
            setLeadtimeDays('');
            setTipoContenitore('');
            setQuantity('');
        }
    }, [kanbanChainId]);


    const handleSubmit = async (e) => {
        e.preventDefault();
        const kanbanData = {
            ...(isEdit ? {
                leadtime_days: parseInt(leadtimeDays, 10),
                tipo_contenitore: tipoContenitore,
                quantity: parseFloat(quantity),
            } : { // Include all fields in create mode
                kanban_chain_id: parseInt(kanbanChainId, 10),
                status_chain_id: parseInt(statusChainId, 10),
                status_current: parseInt(statusCurrent, 10),
                leadtime_days: parseInt(leadtimeDays, 10),
                tipo_contenitore: tipoContenitore,
                quantity: parseFloat(quantity),
                is_active: true // Default to true for new kanbans
            })
        };

        try {
            if (isEdit) {
                await api.put(`/kanbans/${id}`, kanbanData);
            } else {
                await api.post('/kanbans', kanbanData);
            }
            navigate('/kanbans');
        } catch (error) {
            console.error('Error saving kanban:', error);
        }
    };

    return (
        <div>
            <h2>{isEdit ? 'Edit Kanban' : 'Create New Kanban'}</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Kanban Chain:</label>
                    <select value={kanbanChainId} onChange={(e) => setKanbanChainId(e.target.value)} required> {/* Required for Kanban Chain */}
                        <option value="">Select Kanban Chain</option>
                        {availableKanbanChains.map(chain => (
                            <option key={chain.id} value={chain.id}>{chain.id} (Customer: {chain.customer_name}, Product: {chain.product_name})</option>
                        ))}
                    </select>
                </div>
                <div>
                    <label>Status Chain:</label>
                    <input type="text" value={statusChainId} disabled /> {/* Disabled and pre-filled */}
                </div>
                <div>
                    <label>Current Status:</label>
                    <select value={statusCurrent} onChange={(e) => setStatusCurrent(e.target.value)} required> {/* Required for Current Status */}
                        <option value="">Select Status</option>
                        {availableStatuses.map(status => (
                            <option key={status.status_id} value={status.status_id}>{status.name}</option>
                        ))}
                    </select>
                </div>
                <div>
                    <label>Lead Time (days):</label>
                    <input type="number" value={leadtimeDays} disabled /> {/* Disabled and pre-filled */}
                </div>
                <div>
                    <label>Container Type:</label>
                    <input type="text" value={tipoContenitore} disabled /> {/* Disabled and pre-filled */}
                </div>
                <div>
                    <label>Quantity:</label>
                    <input type="number" value={quantity} disabled /> {/* Disabled and pre-filled */}
                </div>

                <div>
                    <label>Is Active:</label>
                    <input type="checkbox" checked={true} disabled={true} />{/* Disabled and default checked to true */}
                </div>


                <button type="submit">{isEdit ? 'Update Kanban' : 'Create Kanban'}</button>
                <Link to="/kanbans">Cancel</Link>
            </form>
        </div>
    );
};

export default KanbanForm;