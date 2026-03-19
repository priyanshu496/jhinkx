"use client";

import { useEffect, useState, useRef } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { Send, ArrowLeft, Loader2, Users, Activity, MessageSquare } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";

export default function SpaceChat() {
  const params = useParams();
  const spaceId = params.id as string;
  const router = useRouter();

  const [user, setUser] = useState<any>(null);
  const [messages, setMessages] = useState<any[]>([]);
  const [newMessage, setNewMessage] = useState("");
  const [isConnected, setIsConnected] = useState(false);
  
  const wsRef = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/signin");
      return;
    }

    const initChat = async () => {
      try {
        const headers = { Authorization: `Bearer ${token}` };

        // A. Get the user's profile
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/users/me`, { headers });
        if (!res.ok) throw new Error("Failed to authenticate");
        const rawData = await res.json();
        
        // DEBUG: Print exactly what Go is sending us so we can see it!
        console.log("Raw Profile Data from Go:", rawData);

        // BULLETPROOF ID EXTRACTOR
        // Search through the object to find the user data, regardless of how Go wrapped it
        const actualUser = rawData.user || rawData.data || rawData.User || rawData;
        
        // Aggressively search for the ID field (accounting for capitalized or lowercase variations)
        const actualUserId = actualUser.id || actualUser.ID || actualUser.userid || actualUser.UserID;

        if (!actualUserId) {
           console.error("CRITICAL ERROR: Could not find User ID. Look at the 'Raw Profile Data' log above.");
           alert("Failed to load user profile correctly. Please check the browser console.");
           return; // STOP EXECUTION! Do not connect to WS with an undefined ID!
        }

        setUser(actualUser);

        // B. Fetch Chat History
        const historyRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/spaces/${spaceId}/messages`, { headers });
        if (historyRes.ok) {
          const historyData = await historyRes.json();
          setMessages(historyData.messages || []);
        }

        // C. Connect to the Go WebSocket Hub safely!
        const wsUrl = `${process.env.NEXT_PUBLIC_API_URL?.replace("http", "ws")}/ws/spaces/${spaceId}?userId=${actualUserId}`;
        const ws = new WebSocket(wsUrl);

        ws.onopen = () => {
          setIsConnected(true);
        };

        ws.onmessage = (event) => {
          const incomingMessage = JSON.parse(event.data);
          setMessages((prev) => [...prev, incomingMessage]);
        };

        ws.onclose = () => {
          setIsConnected(false);
        };

        wsRef.current = ws;
      } catch (err) {
        console.error("Failed to initialize chat", err);
      }
    };

    initChat();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [spaceId, router]);

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim() || !wsRef.current || !isConnected) return;

    const payload = { content: newMessage };
    wsRef.current.send(JSON.stringify(payload));
    setNewMessage(""); 
  };

  if (!user) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-[#f8fafc]">
        <Loader2 className="h-8 w-8 animate-spin text-blue-600 mb-4" />
        <p className="text-slate-500 font-medium">Connecting to secure server...</p>
      </div>
    );
  }

  // Helper function to safely find the user's ID for the chat bubbles
  const getMyId = () => user?.id || user?.ID || user?.userid || user?.UserID;

  return (
    <div className="min-h-screen bg-[#f8fafc] relative font-sans selection:bg-sky-200 flex flex-col">
      <div className="absolute top-0 left-0 w-full h-96 bg-gradient-to-b from-blue-50 to-transparent pointer-events-none"></div>

      <div className="relative z-10 flex-grow max-w-4xl mx-auto w-full flex flex-col p-4 sm:p-6 lg:p-8 h-screen">
        
        <div className="flex justify-between items-center mb-6 bg-white/80 backdrop-blur-md border border-slate-200 p-4 rounded-2xl shadow-sm">
          <div className="flex items-center space-x-4">
            <Link href="/dashboard">
              <Button variant="ghost" size="icon" className="hover:bg-slate-100 rounded-xl">
                <ArrowLeft className="w-5 h-5 text-slate-600" />
              </Button>
            </Link>
            <div>
              <h1 className="text-lg font-bold text-slate-900 flex items-center gap-2">
                Squad Space <span className="text-slate-400 text-sm font-normal">#{spaceId.substring(0, 6)}</span>
              </h1>
              <div className="flex items-center text-xs font-medium text-emerald-600 mt-0.5">
                {isConnected ? (
                  <>
                    <Activity className="w-3.5 h-3.5 mr-1" />
                    Connected & Live
                  </>
                ) : (
                  <span className="text-amber-500 flex items-center">
                    <Loader2 className="w-3 h-3 mr-1 animate-spin" /> Reconnecting...
                  </span>
                )}
              </div>
            </div>
          </div>
          <div className="w-10 h-10 bg-blue-50 rounded-full flex items-center justify-center border border-blue-100">
            <Users className="w-5 h-5 text-blue-600" />
          </div>
        </div>

        <Card className="flex-grow flex flex-col bg-white border-slate-200 shadow-sm rounded-3xl overflow-hidden mb-6">
          <div className="flex-grow p-6 overflow-y-auto space-y-4">
            {messages.length === 0 ? (
              <div className="h-full flex flex-col items-center justify-center text-center space-y-3 opacity-60">
                <MessageSquare className="w-10 h-10 text-slate-300" />
                <p className="text-slate-500 font-medium">No messages yet.</p>
                <p className="text-sm text-slate-400">Be the first to say hello to the squad!</p>
              </div>
            ) : (
              messages.map((msg, index) => {
                const msgUserId = msg.user_id || msg.UserID;
                const msgContent = msg.content || msg.Content;
                
                // Compare the message's user ID with your safely extracted ID
                const isMe = msgUserId === getMyId();
                
                return (
                  <div key={index} className={`flex ${isMe ? "justify-end" : "justify-start"}`}>
                    <div className={`max-w-[75%] rounded-2xl px-5 py-3 ${
                      isMe 
                        ? "bg-blue-600 text-white rounded-br-sm shadow-md" 
                        : "bg-slate-100 text-slate-900 rounded-bl-sm border border-slate-200"
                    }`}>
                      {!isMe && (
                        <p className="text-xs font-bold mb-1 opacity-50">
                          {msgUserId?.substring(0, 6) || "User"}
                        </p>
                      )}
                      <p className="text-sm leading-relaxed">{msgContent}</p>
                    </div>
                  </div>
                );
              })
            )}
            <div ref={messagesEndRef} />
          </div>

          <div className="p-4 bg-slate-50 border-t border-slate-100">
            <form onSubmit={handleSendMessage} className="flex gap-3">
              <Input
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                placeholder="Type a message to the squad..."
                className="flex-grow h-12 bg-white border-slate-200 rounded-xl focus-visible:ring-blue-500"
                disabled={!isConnected}
              />
              <Button 
                type="submit" 
                disabled={!newMessage.trim() || !isConnected}
                className="h-12 w-12 rounded-xl bg-blue-600 hover:bg-blue-700 shadow-sm transition-all p-0 flex items-center justify-center"
              >
                <Send className="w-5 h-5 text-white ml-0.5" />
              </Button>
            </form>
          </div>
        </Card>

      </div>
    </div>
  );
}