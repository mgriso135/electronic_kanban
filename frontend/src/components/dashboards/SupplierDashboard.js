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

    const fetchDashboardData = useCallback(async () => {
        if (!selectedSupplierId) return;
        setIsLoading(true);
        setIsDashboardDataReady(false);
        try {
            const kanbanResponse = await api.get(`/dashboards/supplier/${selectedSupplierId}`);
            let kanbans = kanbanResponse.data.kanbans_by_product;
            kanbans = sortKanbansByCustomerSupplierAndDate(kanbans);
            setKanbansByProduct(kanbans);
        } catch (error) {
            console.error("Error fetching data for supplier dashboard", error);
        } finally {
            setIsLoading(false);
            setIsDashboardDataReady(true);
        }
    }, [selectedSupplierId, sortKanbansByCustomerSupplierAndDate]);

    useEffect(() => {
        const fetchSuppliers = async () => {
            try {
                const accountsResponse = await api.get('/accounts');
                const suppliers = accountsResponse.data;
                setAvailableSuppliers(suppliers);

                const initialSupplierId = supplierId || (suppliers.length > 0 ? suppliers[0].id : '');
                setSelectedSupplierId(initialSupplierId);
            } catch (error) {
                console.error("Error fetching suppliers", error);
            }
        };
        fetchSuppliers();
    }, [supplierId]);

    useEffect(() => {
        fetchDashboardData();

        const intervalId = setInterval(() => {
            console.log("SupplierDashboard: Polling dashboard data...");
            fetchDashboardData();
        }, 5000);

        return () => clearInterval(intervalId);
    }, [fetchDashboardData]);

    const handleSupplierChange = (e) => {
        const newSupplierId = e.target.value;
        setSelectedSupplierId(newSupplierId);

        if (newSupplierId) {
            navigate(`/supplier-dashboard/${newSupplierId}`);
        } else {
            navigate(`/supplier-dashboard`);
        }
    };

    const handleKanbanUpdate = useCallback((updatedKanban, productID) => {
        setKanbansByProduct(prevKanbansByProduct => {
            const updatedKanbansByProduct = { ...prevKanbansByProduct };
            for (const product in updatedKanbansByProduct) {
                updatedKanbansByProduct[product] = updatedKanbansByProduct[product].map(k => {
                    if (k.kanban_id === updatedKanban.kanban_id) {
                        return { ...k, status_name: updatedKanban.status_name, status_color: updatedKanban.status_color, customer_supplier: updatedKanban.customer_supplier, status_current: updatedKanban.status_current }; // UPDATED - Removed customer/supplier name updates
                    } else {
                        return { ...k, 
                            status_name: k.status_name,
                            status_color: k.status_color,
                            customer_supplier: k.customer_supplier,
                            status_current: k.status_current,
                            customer_name: k.customer_name, // Keep customer_name preservation (SupplierDashboard)
                            supplier_name: k.supplier_name, // Keep supplier_name preservation (CustomerDashboard)
                            ...k 
                         }; 
                    }
                });
            }
            return updatedKanbansByProduct;
        });
    }, []);

    const handleStatusChangeSuccess = useCallback(() => {
        console.log("SupplierDashboard - handleStatusChangeSuccess CALLED - Re-fetching Kanban data");
        fetchDashboardData(); // Re-fetch dashboard data
    }, [fetchDashboardData]);

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