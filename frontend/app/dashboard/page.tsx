"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { LogOut, Users, MessageSquare, Loader2, Target, Plus } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function Dashboard() {
  const [user, setUser] = useState<any>(null);
  const [spaces, setSpaces] = useState<any[]>([]); // Safely initialized as an empty array
  const [isSearching, setIsSearching] = useState(false);
  const [groupSize, setGroupSize] = useState(4);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/signin");
      return;
    }

    const fetchDashboardData = async () => {
      try {
        const headers = { Authorization: `Bearer ${token}` };
        
        // 1. Fetch User
        const userRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/users/me`, { headers });
        if (userRes.ok) setUser(await userRes.json());

        // 2. Fetch Spaces (WITH SAFETY CHECK)
        const spacesRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/spaces`, { headers });
        if (spacesRes.ok) {
          const data = await spacesRes.json();
          // Safety: Force it to be an array, even if the API returns null or an object
          if (Array.isArray(data)) {
            setSpaces(data);
          } else if (data && Array.isArray(data.spaces)) {
            setSpaces(data.spaces);
          } else {
            setSpaces([]);
          }
        }
      } catch (err) {
        console.error("Failed to load dashboard data", err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchDashboardData();
  }, [router]);

  const handleFindMatch = async () => {
    setIsSearching(true);
    const token = localStorage.getItem("token");
    
    // Capture the current space IDs before we start searching
    const existingIds = new Set(spaces.map(s => s.id));

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/spaces/match`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ preferred_group_size: groupSize }),
      });

      if (res.ok) {
        const interval = setInterval(async () => {
          const checkRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/spaces`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          
          if (checkRes.ok) {
            const currentData = await checkRes.json();
            
            // THE FIX: Correctly extract the array from the JSON object
            let currentSpaces = [];
            if (currentData && Array.isArray(currentData.spaces)) {
              currentSpaces = currentData.spaces;
            } else if (Array.isArray(currentData)) {
              currentSpaces = currentData;
            }
            
            // SMART CHECK: Is there a space ID that wasn't in our 'existingIds' set?
            const newMatch = currentSpaces.find((s: any) => !existingIds.has(s.id));

            if (newMatch) {
              clearInterval(interval);
              setIsSearching(false);
              setSpaces(currentSpaces);


              // 1. TRIGGER THE BEAUTIFUL SUCCESS TOAST!
              toast.success("Match Found! 🎉", {
                description: "Squad assembled. Deploying to chat room...",
                duration: 3000,
              });

              // 2. WAIT 1.5 SECONDS SO THEY CAN READ IT, THEN TELEPORT!
              setTimeout(() => {
                router.push(`/spaces/${newMatch.id}`);
              }, 1500);
            }
          }
        }, 2500); // Poll every 2.5 seconds
      } else {
        const errorData = await res.json().catch(() => ({}));
        alert(`Error: ${errorData.error || "Matchmaking failed"}`);
        setIsSearching(false);
      }
    } catch (err) {
      console.error(err);
      setIsSearching(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    router.push("/signin");
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-[#f8fafc]">
        <Loader2 className="h-8 w-8 animate-spin text-slate-400 mb-4" />
        <p className="text-slate-500 font-medium animate-pulse">Loading your dashboard...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#f8fafc] relative font-sans selection:bg-sky-200 pb-20">
      {/* Soft Background Glows */}
      <div className="absolute top-0 left-1/4 w-96 h-96 bg-blue-100 rounded-full mix-blend-multiply filter blur-3xl opacity-50 pointer-events-none"></div>
      <div className="absolute top-40 right-1/4 w-96 h-96 bg-sky-100 rounded-full mix-blend-multiply filter blur-3xl opacity-50 pointer-events-none"></div>

      <div className="relative z-10 max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 pt-12">
        
        {/* Header Bar */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-10 gap-4">
          <div>
            <h1 className="text-3xl font-bold tracking-tight text-slate-900">
              Welcome back, <span className="text-blue-600">{user?.username}</span>
            </h1>
            <p className="text-slate-500 mt-1">Ready to deploy? Find a squad or jump back into a chat.</p>
          </div>
          <Button variant="outline" onClick={handleLogout} className="bg-white/50 backdrop-blur-sm border-slate-200 hover:bg-red-50 hover:text-red-600 hover:border-red-200 transition-colors rounded-xl h-11 px-6">
            <LogOut className="w-4 h-4 mr-2" />
            Sign Out
          </Button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          
          {/* Matchmaker Panel (Left Column) */}
          <div className="lg:col-span-5">
            <Card className="bg-white/80 backdrop-blur-xl border-white/40 shadow-xl rounded-3xl overflow-hidden relative">
              {isSearching && (
                <div className="absolute top-0 left-0 w-full h-1 bg-slate-100">
                  <div className="h-full bg-blue-600 animate-[progress_2s_ease-in-out_infinite] w-1/2 rounded-full"></div>
                </div>
              )}
              
              <CardHeader className="pt-8 pb-4">
                <div className="w-12 h-12 bg-blue-50 rounded-2xl flex items-center justify-center mb-4 border border-blue-100/50">
                  <Target className="w-6 h-6 text-blue-600" />
                </div>
                <CardTitle className="text-xl font-bold text-slate-900">Find a Squad</CardTitle>
                <CardDescription>Select your team size to enter the Kafka matchmaking queue.</CardDescription>
              </CardHeader>

              <CardContent className="pb-8">
                {!isSearching ? (
                  <div className="space-y-6">
                    <div className="space-y-3">
                      <label className="text-sm font-semibold text-slate-700">Required Players</label>
                      <select
                        value={groupSize}
                        onChange={(e) => setGroupSize(Number(e.target.value))}
                        className="w-full h-12 px-4 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all appearance-none font-medium text-slate-700"
                      >
                        {[2, 3, 4, 5, 6].map(size => (
                          <option key={size} value={size}>Group of {size}</option>
                        ))}
                      </select>
                    </div>
                    <Button 
                      onClick={handleFindMatch}
                      className="w-full h-12 bg-slate-900 hover:bg-slate-800 text-white rounded-xl shadow-md transition-all hover:-translate-y-0.5"
                    >
                      <Users className="w-4 h-4 mr-2" />
                      Start Matchmaking
                    </Button>
                  </div>
                ) : (
                  <div className="py-8 flex flex-col items-center justify-center text-center space-y-4">
                    <div className="relative w-16 h-16 flex items-center justify-center">
                      <div className="absolute inset-0 rounded-full border-4 border-blue-100"></div>
                      <div className="absolute inset-0 rounded-full border-4 border-blue-600 border-t-transparent animate-spin"></div>
                      <Users className="w-6 h-6 text-blue-600 animate-pulse" />
                    </div>
                    <div>
                      <h3 className="font-bold text-slate-900">Searching for players...</h3>
                      <p className="text-sm text-slate-500 mt-1">Listening to Kafka events</p>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Active Spaces Panel (Right Column) */}
          <div className="lg:col-span-7">
            <Card className="bg-white/80 backdrop-blur-xl border-white/40 shadow-xl rounded-3xl h-full">
              <CardHeader className="pt-8 pb-4">
                <div className="flex justify-between items-center">
                  <div>
                    <CardTitle className="text-xl font-bold text-slate-900">Active Spaces</CardTitle>
                    <CardDescription>Your live WebSocket chat rooms.</CardDescription>
                  </div>
                  <div className="bg-blue-50 text-blue-700 text-xs font-bold px-3 py-1.5 rounded-full border border-blue-100">
                    {spaces.length} Online
                  </div>
                </div>
              </CardHeader>
              
              <CardContent>
                {spaces.length === 0 ? (
                  <div className="py-12 flex flex-col items-center justify-center text-center border-2 border-dashed border-slate-200 rounded-2xl bg-slate-50/50">
                    <div className="w-12 h-12 bg-white rounded-full flex items-center justify-center shadow-sm mb-4">
                      <MessageSquare className="w-5 h-5 text-slate-400" />
                    </div>
                    <h3 className="font-semibold text-slate-900">No active spaces</h3>
                    <p className="text-sm text-slate-500 max-w-[250px] mt-1">
                      You aren't in any chat rooms right now. Use the matchmaker to find a squad!
                    </p>
                  </div>
                ) : (
                  <div className="grid gap-4">
                    {spaces.map((space: any) => (
                      <Link 
                        href={`/spaces/${space.id}`} 
                        key={space.id}
                        className="group flex items-center justify-between p-4 sm:p-5 bg-white border border-slate-200 rounded-2xl hover:border-blue-300 hover:shadow-md transition-all hover:bg-blue-50/30"
                      >
                        <div className="flex items-center space-x-4">
                          <div className="w-10 h-10 bg-slate-100 rounded-full flex items-center justify-center group-hover:bg-blue-100 transition-colors">
                            <MessageSquare className="w-5 h-5 text-slate-500 group-hover:text-blue-600 transition-colors" />
                          </div>
                          <div>
                            <p className="font-bold text-slate-900 group-hover:text-blue-900 transition-colors">
                              Space #{space.id.substring(0, 6)}
                            </p>
                            <p className="text-sm text-slate-500">
                              Target Size: {space.target_size} Players
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center">
                          <span className="flex items-center text-xs font-semibold text-emerald-600 bg-emerald-50 px-2.5 py-1 rounded-full border border-emerald-100">
                            <span className="w-1.5 h-1.5 bg-emerald-500 rounded-full mr-1.5 animate-pulse"></span>
                            Live
                          </span>
                        </div>
                      </Link>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

        </div>
      </div>

      {/* Add a custom keyframe animation for the searching progress bar */}
      <style dangerouslySetInnerHTML={{__html: `
        @keyframes progress {
          0% { transform: translateX(-100%); }
          100% { transform: translateX(200%); }
        }
      `}} />
    </div>
  );
}