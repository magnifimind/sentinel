import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Layout } from "./components/Layout";
import { Activity } from "./pages/Activity";
import { AuditTrail } from "./pages/AuditTrail";
import { Agents } from "./pages/Agents";
import { Policies } from "./pages/Policies";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route index element={<Activity />} />
          <Route path="audit" element={<AuditTrail />} />
          <Route path="agents" element={<Agents />} />
          <Route path="policies" element={<Policies />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
