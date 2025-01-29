import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../services/api';
import KanbanCard from '../kanbans/KanbanCard';

const SupplierDashboard = () => {
    const { supplierId } = useParams();
    const [kanbansByProduct, setKanbansByProduct] = useState({});
    const [availableSuppliers, setAvailableSuppliers] = useState([]);
    const [selectedSupplierId, setSelectedSupplierId] = useState(supplierId || '');


    useEffect(() => {
        const fetchSuppliersAndKanbans = async () => {
            try {
                const accountsResponse = await api.get('/accounts');
                // Filter for suppliers only
                const suppliers = accountsResponse.data;
                setAvailableSuppliers(suppliers);


                if(selectedSupplierId){
                    const kanbanResponse = await api.get(`/dashboards/supplier/${selectedSupplierId}`);
                   let kanbans = kanbanResponse.data.kanbans_by_product;


                    for (const product in kanbans) {
                        kanbans[product].sort((a, b) => {
                          const aMatch = a.customer_supplier === 1 ? -1 : 1; // Move matching to front
                          const bMatch = b.customer_supplier === 1 ? -1 : 1; // Move matching to front
                         if (aMatch !== bMatch) {
                            return aMatch - bMatch;
                           }
                         return new Date(a.data_aggiornamento) - new Date(b.data_aggiornamento); // Sort by date
                        });
                     }

                    setKanbansByProduct(kanbans);
                }

            } catch (error) {
                console.error("Error fetching data for supplier dashboard", error);
            }
        };

        fetchSuppliersAndKanbans();
    }, [selectedSupplierId]);

  const handleSupplierChange = (e) => {
    setSelectedSupplierId(e.target.value);
  };


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

            {Object.keys(kanbansByProduct).length > 0 ? (
                Object.keys(kanbansByProduct).map(product => (
                    <div key={product}>
                        <h3>{product}</h3>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '20px'}}>
                        {kanbansByProduct[product].map(kanban => {
                            console.log("SupplierDashboard: Kanban Object:", kanban); // DEBUG: Log kanban object
                            return (
                                <KanbanCard key={kanban.kanban_id} kanban={kanban} dashboardType="supplier" setKanbans={setKanbansByProduct} />
                            )
                        })}
                        </div>
                    </div>
                ))
            ) : (
                <p>Please select a supplier to view Kanban data</p>
             )}

        </div>
    );
};

export default SupplierDashboard;