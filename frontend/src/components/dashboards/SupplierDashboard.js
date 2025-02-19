import React, { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../services/api';
import KanbanCard from '../kanbans/KanbanCard';

const SupplierDashboard = () => {
    const { supplierId } = useParams();
    const navigate = useNavigate(); 
    const [kanbansByProduct, setKanbansByProduct] = useState({});
    const [availableSuppliers, setAvailableSuppliers] = useState([]);
    const [selectedSupplierId, setSelectedSupplierId] = useState(supplierId || '');
    const [isLoading, setIsLoading] = useState(false);
    const [isDashboardDataReady, setIsDashboardDataReady] = useState(false); // ADD NEW isDashboardDataReady STATE


    useEffect(() => {
        fetchSuppliersAndKanbans();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedSupplierId]);

    const fetchSuppliersAndKanbans = async () => {
        setIsLoading(true);
        setIsDashboardDataReady(false); // **SET isDashboardDataReady to FALSE at start of fetch**
        try {
            const accountsResponse = await api.get('/accounts');
            // Filter for suppliers only
            const suppliers = accountsResponse.data;
            setAvailableSuppliers(suppliers);

            const initialSupplierId = supplierId || (suppliers.length > 0 ? suppliers[0].id : '');
            setSelectedSupplierId(initialSupplierId);


            if (initialSupplierId) {
                const kanbanResponse = await api.get(`/dashboards/supplier/${initialSupplierId}`);
                let kanbans = kanbanResponse.data.kanbans_by_product;
                kanbans = sortKanbansByCustomerSupplierAndDate(kanbans);
                setKanbansByProduct(kanbans);
            } else {
                setKanbansByProduct({});
            }


        } catch (error) {
            console.error("Error fetching data for supplier dashboard", error);
        } finally {
            setIsLoading(false);
            setIsDashboardDataReady(true); // **SET isDashboardDataReady to TRUE in finally block** - Data loading is complete
        }
    };

    const handleSupplierChange = (e) => {
        const newSupplierId = e.target.value;
        setSelectedSupplierId(newSupplierId);

        if (newSupplierId) {
            navigate(`/supplier-dashboard/${newSupplierId}`);
        } else {
            navigate(`/supplier-dashboard`);
        }
    };


    const sortKanbansByCustomerSupplierAndDate = useCallback((kanbans) => {
        const sortedKanbans = { ...kanbans }; // Create a copy to avoid modifying original
        for (const product in sortedKanbans) {
            sortedKanbans[product].sort((a, b) => {
                const aMatch = a.customer_supplier === 1 ? -1 : 1; // Move matching to front
                const bMatch = b.customer_supplier === 1 ? -1 : 1; // Move matching to front
                if (aMatch !== bMatch) {
                    return aMatch - bMatch;
                }
                return new Date(a.data_aggiornamento) - new Date(b.data_aggiornamento); // Sort by date
            });
        }
        return sortedKanbans;
    }, []);

    const handleKanbanUpdate = useCallback((updatedKanban, productID) => {
        setKanbansByProduct(prevKanbansByProduct => {
            const updatedKanbansByProduct = { ...prevKanbansByProduct };
            for (const product in updatedKanbansByProduct) {
                updatedKanbansByProduct[product] = updatedKanbansByProduct[product].map(k => {
                    if (k.kanban_id === updatedKanban.kanban_id) {
                        return { ...k, status_name: updatedKanban.status_name, status_color: updatedKanban.status_color, customer_supplier: updatedKanban.customer_supplier, status_current: updatedKanban.status_current, supplier_name: updatedKanban.supplier_name };
                    } else {
                        return { ...k, 
                            status_name: k.status_name, // **PRESERVE status_name**
                            status_color: k.status_color, // **PRESERVE status_color**
                            customer_supplier: k.customer_supplier, // **PRESERVE customer_supplier**
                            status_current: k.status_current, // **PRESERVE status_current**
                            supplier_name: k.supplier_name, // **PRESERVE supplier_name**
                            ...k // **IMPORTANT: ALSO PRESERVE OTHER EXISTING PROPERTIES using spread operator**
                         }; // Return k and PRESERVE ALL EXISTING PROPERTIES for other kanbans
                    }
                });
            }
            return updatedKanbansByProduct;
        });
    }, []);

    const handleStatusChangeSuccess = useCallback(() => {
        console.log("SupplierDashboard - handleStatusChangeSuccess CALLED - Re-fetching Kanban data");
        fetchSuppliersAndKanbans(); // Re-fetch data from API
    }, [fetchSuppliersAndKanbans]);

    return (
        <div>
            <h2>Supplier Dashboard</h2>

            <div>
                <label>Select Supplier:</label>
                <select value={selectedSupplierId} onChange={handleSupplierChange} disabled={isLoading}>
                    <option value="">Select a Supplier</option>
                    {availableSuppliers.map(supplier => (
                        <option key={supplier.id} value={supplier.id}>{supplier.name}</option>
                    ))}
                </select>
            </div>

            {isLoading ? (
                <p>Loading Kanban data...</p>
            ) : (
                Object.keys(kanbansByProduct).length > 0 ? (
                    Object.keys(kanbansByProduct).map(product => (
                        <div key={product}>
                            <h3>{product}</h3>
                            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '20px' }}>
                            {kanbansByProduct[product].map(kanban => (
                            <KanbanCard
                                key={kanban.kanban_id}
                                kanban={kanban}
                                dashboardType="supplier"
                                setKanbans={handleKanbanUpdate}
                                productID={product}
                                isDashboardDataReady={isDashboardDataReady}
                                onStatusChangeSuccess={handleStatusChangeSuccess} 
                            />
                        ))}
                            </div>
                        </div>
                    ))
                ) : (
                    <p>Please select a supplier to view Kanban data</p>
                ))}

            </div>
        );
    };

export default SupplierDashboard;