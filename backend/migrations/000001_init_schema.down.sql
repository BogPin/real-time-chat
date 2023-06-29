DROP TABLE public.users;
DROP SEQUENCE public.users_id_seq;
DROP TABLE public.participants;
DROP TABLE public.messages;
DROP SEQUENCE public.messages_id_seq;
DROP TABLE public.chats;
DROP SEQUENCE public.chats_id_seq;

DELETE TYPE public.message_type;
DELETE TYPE public.participant_role;
