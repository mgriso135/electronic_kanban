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
    const [initialNoOfActiveKanbans, setInitialNoOfActiveKanbans] = useState(1);
    const [isEdit, setIsEdit] = useState(false);
    const [availableStatusChains, setAvailableStatusChains] = useState([]);
    const [availableCustomers, setAvailableCustomers] = useState([]);
    const [availableSuppliers, setAvailableSuppliers] = useState([]);
    const [availableProducts, setAvailableProducts] = useState([]);
    const [showDeleteKanbanMessage, setShowDeleteKanbanMessage] = useState(false);

    useEffect(() => {
        const fetchDropdownData = async () => {
            try {
                const statusChainsResponse = await api.get('/status-chains');
                setAvailableStatusChains(statusChainsResponse.data);
                const customersResponse = await api.get('/accounts');
                setAvailableCustomers(customersResponse.data);
                const suppliersResponse = await api.get('/accounts');
                setAvailableSuppliers(suppliersResponse.data);
                const productsResponse = await api.get('/products');
                setAvailableProducts(productsResponse.data);

            } catch (error) {
                console.error("Error fetching dropdown data:", error);
            }
        };
        fetchDropdownData();

        if (id) {
            setIsEdit(true);
            const fetchKanbanChain = async () => {
                try {
                    const response = await api.get(`/kanban-chains/${id}`);
                    const kcData = response.data;
                    setClienteId(String(kcData.cliente_id));
                    setProdottoCodice(kcData.prodotto_codice);
                    setFornitoreId(String(kcData.fornitore_id));
                    setLeadtimeDays(String(kcData.leadtime_days));
                    setQuantity(String(kcData.quantity));
                    setTipoContenitore(kcData.tipo_contenitore);
                    setStatusChainId(String(kcData.status_chain_id));
                    setNoOfActiveKanbans(String(kcData.no_of_active_kanbans));
                    setInitialNoOfActiveKanbans(kcData.no_of_active_kanbans);
                } catch (error) {
                    console.error('Error fetching kanban chain:', error);
                }
            };
            fetchKanbanChain();
        } else {
            setInitialNoOfActiveKanbans(0);
        }
    }, [id]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setShowDeleteKanbanMessage(false);

        const kanbanChainData = {
            kanban_chain: {
                id: parseInt(id, 10) || 0, // Include ID for edit, default to 0 for create
                cliente_id: parseInt(clienteId, 10),
                prodotto_codice: prodottoCodice,
                fornitore_id: parseInt(fornitoreId, 10),
                leadtime_days: parseInt(leadtimeDays, 10),
                quantity: parseFloat(quantity),
                tipo_contenitore: tipoContenitore,
                status_chain_id: parseInt(statusChainId, 10),
                no_of_active_kanbans: parseInt(noOfActiveKanbans, 10) // Include in payload
            },
             no_of_initial_kanbans: parseInt(noOfActiveKanbans, 10)
        };


        const newNoOfActiveKanbans = parseInt(noOfActiveKanbans, 10);
        if (isEdit && newNoOfActiveKanbans > initialNoOfActiveKanbans) {
             kanbanChainData.no_of_initial_kanbans = newNoOfActiveKanbans - initialNoOfActiveKanbans;
         } else if (isEdit && newNoOfActiveKanbans < initialNoOfActiveKanbans) {
             setShowDeleteKanbanMessage(true);
            return; // Early return to prevent update and show message
        }


        try {
            if (isEdit) {
                 await api.put(`/kanban-chains/${id}`, kanbanChainData);
            } else {
                await api.post('/kanban-chains', kanbanChainData);
            }
            navigate('/kanban-chains');
        } catch (error) {
            console.error('Error saving kanban chain:', error);
        }
    };


    return (
        <div>
            <h2>{isEdit ? 'Edit Kanban Chain' : 'Create New Kanban Chain'}</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Customer:</label>
                    <select value={clienteId} onChange={(e) => setClienteId(e.target.value)} disabled={isEdit}>
                        <option value="">Select Customer</option>
                        {availableCustomers.map(customer => (
                            <option key={customer.id} value={customer.id}>{customer.name}</option>
                        ))}
                    </select>
                </div>
                <div>
                    <label>Supplier:</label>
                    <select value={fornitoreId} onChange={(e) => setFornitoreId(e.target.value)} disabled={isEdit}>
                        <option value="">Select Supplier</option>
                        {availableSuppliers.map(supplier => (
                            <option key={supplier.id} value={supplier.id}>{supplier.name}</option>
                        ))}
                    </select>
                </div>
                <div>
                    <label>Product:</label>
                    <select value={prodottoCodice} onChange={(e) => setProdottoCodice(e.target.value)} disabled={isEdit}>
                        <option value="">Select Product</option>
                        {availableProducts.map(product => (
                            <option key={product.product_id} value={product.product_id}>{product.name} ({product.product_id})</option>
                        ))}
                    </select>
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
                    <select value={statusChainId} onChange={(e) => setStatusChainId(e.target.value)} disabled={isEdit}>
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
                {showDeleteKanbanMessage && (
                    <p>To remove active kanbans, please delete them manually from the Kanban list.</p>
                )}

                <button type="submit">{isEdit ? 'Update Kanban Chain' : 'Create Kanban Chain'}</button>
                <Link to="/kanban-chains">Cancel</Link>
            </form>
        </div>
    );
};

export default KanbanChainForm;