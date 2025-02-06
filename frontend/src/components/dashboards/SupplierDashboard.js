import React, { useState, useEffect, useCallback } from 'react'; // Import useCallback
import { useParams, useNavigate, useLocation } from 'react-router-dom'; // Import useNavigate, useLocation
import api from '../../services/api';
import KanbanCard from '../kanbans/KanbanCard';

const SupplierDashboard = () => {
    const { supplierId } = useParams();
    const navigate = useNavigate(); // Hook for navigation
    const location = useLocation(); // Hook to get current location
    const [kanbansByProduct, setKanbansByProduct] = useState({});
    const [availableSuppliers, setAvailableSuppliers] = useState([]);
    const [selectedSupplierId, setSelectedSupplierId] = useState(supplierId || '');

    useEffect(() => {
        fetchSuppliersAndKanbans();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedSupplierId]);

    const fetchSuppliersAndKanbans = async () => {
        try {
            const accountsResponse = await api.get('/accounts');
            const suppliers = accountsResponse.data;
            setAvailableSuppliers(suppliers);

            // If supplierId is in URL, use it, otherwise default to first supplier or empty selection
            const initialSupplierId = supplierId || (suppliers.length > 0 ? suppliers[0].id : '');
            setSelectedSupplierId(initialSupplierId);

            if (initialSupplierId) { // Fetch kanbans only if a supplier is selected
                const kanbanResponse = await api.get(`/dashboards/supplier/${initialSupplierId}`);
                let kanbans = kanbanResponse.data.kanbans_by_product;
                kanbans = sortKanbansByCustomerSupplierAndDate(kanbans);
                setKanbansByProduct(kanbans);
            } else {
                setKanbansByProduct({}); // Clear kanbans if no supplier selected
            }


        } catch (error) {
            console.error("Error fetching data for supplier dashboard", error);
        }
    };

    const handleSupplierChange = (e) => {
        const newSupplierId = e.target.value;
        setSelectedSupplierId(newSupplierId);

        if (newSupplierId) {
            navigate(`/supplier-dashboard/${newSupplierId}`); // Update URL to include supplierId
        } else {
            navigate(`/supplier-dashboard`); // Navigate to generic dashboard if no supplier selected
        }
    };


    const sortKanbansByCustomerSupplierAndDate = useCallback((kanbans) => { // Use useCallback
        const sortedKanbans = { ...kanbans };
        for (const product in sortedKanbans) {
            sortedKanbans[product].sort((a, b) => {
                const aMatch = a.customer_supplier === 1 ? -1 : 1;
                const bMatch = b.customer_supplier === 1 ? -1 : 1;
                if (aMatch !== bMatch) {
                    return aMatch - bMatch;
                }
                return new Date(a.data_aggiornamento) - new Date(b.data_aggiornamento);
            });
        }
        return sortedKanbans;
    }, []);

    const handleKanbanUpdate = useCallback((updatedKanban, productID) => { // Use useCallback and accept productID
        setKanbansByProduct(prevKanbansByProduct => {
            // Safety checks (already present - keep them)
            if (!prevKanbansByProduct || !prevKanbansByProduct[productID] || !Array.isArray(prevKanbansByProduct[productID])) {
                console.warn("handleKanbanUpdate: prevKanbansByProduct or product data is not properly initialized yet.");
                return prevKanbansByProduct || {};
            }

            const updatedProductKanbans = prevKanbansByProduct[productID].map(k =>
                k.kanban_id === updatedKanban.kanban_id ? updatedKanban : k
            );

            // Create a NEW kanbansByProduct object instead of modifying in-place
            const nextKanbansByProduct = { ...prevKanbansByProduct }; // Start with a copy
            nextKanbansByProduct[productID] = sortKanbansByCustomerSupplierAndDate({ [productID]: updatedProductKanbans })[productID]; // Re-sort and assign

            return nextKanbansByProduct; // Return the completely new object
        });
    }, [sortKanbansByCustomerSupplierAndDate]); // Dependency array includes sort function

    return (
        <div>
            <h2>Supplier Dashboard</h2>

            <div>
                <label>Select Supplier:</label>
                <select value={selectedSupplierId} onChange={handleSupplierChange}>
                    <option value="">Select a Supplier</option>
                    {availableSuppliers.map(supplier => (
                        <option key={supplier.id} value={supplier.id}>{supplier.name}</option>
                    ))}
                </select>
            </div>

            {selectedSupplierId && Object.keys(kanbansByProduct).length > 0 ? ( // Conditionally render kanban data
                Object.keys(kanbansByProduct).map(product => (
                    <div key={product}>
                        <h3>{product}</h3>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '20px' }}>
                            {kanbansByProduct[product].map(kanban => (
                                <KanbanCard
                                    key={kanban.kanban_id}
                                    kanban={kanban}
                                    dashboardType="supplier"
                                    setKanbans={handleKanbanUpdate} // Pass handleKanbanUpdate
                                    productID={product}
                                />
                            ))}
                        </div>
                    </div>
                ))
            ) : (
                selectedSupplierId ? <p>No Kanban data found for this supplier.</p> : <p>Please select a supplier to view Kanban data</p>
            )}

        </div>
    );
};

export default SupplierDashboard;