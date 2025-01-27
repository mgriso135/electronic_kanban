import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../services/api';
import KanbanCard from '../kanbans/KanbanCard';

const SupplierDashboard = () => {
    const { supplierId } = useParams();
  const [kanbansByProduct, setKanbansByProduct] = useState({});


    useEffect(() => {
      const fetchKanbans = async () => {
        try {
            const response = await api.get(`/dashboards/supplier/${supplierId}`);
          setKanbansByProduct(response.data.kanbans_by_product);
        } catch (error) {
            console.error("Error fetching supplier dashboard data", error);
        }
      };
        fetchKanbans();
    }, [supplierId]);

    return (
      <div>
        <h2>Supplier Dashboard</h2>
        {Object.keys(kanbansByProduct).map(product => (
          <div key={product}>
            <h3>{product}</h3>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '20px'}}>
              {kanbansByProduct[product].map(kanban => (
                <KanbanCard key={kanban.kanban_id} kanban={kanban} dashboardType="supplier" setKanbans={setKanbansByProduct} />
                ))}
            </div>
          </div>
        ))}
    </div>
    );
};

export default SupplierDashboard;