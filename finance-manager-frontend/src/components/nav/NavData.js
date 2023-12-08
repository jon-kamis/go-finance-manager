import HomeIcon from '@mui/icons-material/Home';
import TravelExploreIcon from '@mui/icons-material/TravelExplore'
import BarChartIcon from '@mui/icons-material/BarChart';
import SettingsIcon from '@mui/icons-material/Settings';

export const NavData = [
    {
        id: 0,
        text: "Home",
        icon: <HomeIcon/>,
        link: "/"
    },
    {
        id: 1,
        text: "About",
        icon: <TravelExploreIcon/>,
        link: "/about"
    },
    {
        id: 2,
        text: "Users",
        icon: <SettingsIcon/>,
        link: "/users"
    }
]